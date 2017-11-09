package agentd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/sensu/sensu-go/backend/messaging"
	"github.com/sensu/sensu-go/backend/store"
	"github.com/sensu/sensu-go/handler"
	"github.com/sensu/sensu-go/transport"
	"github.com/sensu/sensu-go/types"
)

// A Session is a server-side connection between a Sensu backend server and
// the Sensu agent process via the Sensu transport. It is responsible for
// relaying messages to the message bus on behalf of the agent and from the
// bus to the agent from other daemons. It handles transport handshaking and
// transport channel multiplexing/demultiplexing.
type Session struct {
	ID string

	cfg          SessionConfig
	conn         transport.Transport
	store        store.Store
	handler      *handler.MessageHandler
	stopping     chan struct{}
	wg           *sync.WaitGroup
	sendq        chan *transport.Message
	checkChannel chan interface{}
	bus          messaging.MessageBus
	logger       logrus.FieldLogger
}

func newSessionHandler(s *Session) *handler.MessageHandler {
	handler := handler.NewMessageHandler()
	handler.AddHandler(transport.MessageTypeKeepalive, s.handleKeepalive)
	handler.AddHandler(transport.MessageTypeEvent, s.handleEvent)

	return handler
}

// A SessionConfig contains all of the ncessary information to intiialize
// an agent session.
type SessionConfig struct {
	Organization  string
	Environment   string
	AgentID       string
	User          string
	Subscriptions []string
}

// NewSession creates a new Session object given the triple of a transport
// connection, message bus, and store.
func NewSession(cfg SessionConfig, conn transport.Transport, bus messaging.MessageBus, store store.Store) (*Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	// Validate the agent organization and environment
	ctx := context.TODO()
	if _, err := store.GetEnvironment(ctx, cfg.Organization, cfg.Environment); err != nil {
		return nil, fmt.Errorf("the environment '%s:%s' is invalid", cfg.Organization, cfg.Environment)
	}

	logEntry := logger.WithField("agentId", cfg.AgentID)
	logEntry.WithField("subscriptions", cfg.Subscriptions).Info("agent connected")

	s := &Session{
		ID:           id.String(),
		conn:         conn,
		cfg:          cfg,
		stopping:     make(chan struct{}, 1),
		wg:           &sync.WaitGroup{},
		sendq:        make(chan *transport.Message, 10),
		checkChannel: make(chan interface{}, 100),
		logger:       logEntry,

		store: store,
		bus:   bus,
	}
	s.handler = newSessionHandler(s)
	return s, nil
}

func (s *Session) receiveMessages(out chan *transport.Message) {
	defer close(out)
	for {
		m, err := s.conn.Receive()
		if err != nil {
			switch err := err.(type) {
			case transport.ConnectionError, transport.ClosedError:
				s.logger.Error("recv error: ", err.Error())
				return
			default:
				s.logger.Error("recv error: ", err.Error())
				continue
			}
		}
		out <- m
	}
}

func (s *Session) recvPump() {
	defer func() {
		s.logger.Info("session disconnected - stopping recvPump")
		s.wg.Done()
	}()

	msgChannel := make(chan *transport.Message)
	go s.receiveMessages(msgChannel)
	for {
		select {
		case msg, ok := <-msgChannel:
			if !ok {
				return
			}

			s.logger.Debugf("session - received message: %s", string(msg.Payload))
			err := s.handler.Handle(msg.Type, msg.Payload)
			if err != nil {
				s.logger.WithError(err).WithField("msg", msg).Error("error handling message")
			}
		case <-s.stopping:
			return
		}
	}
}

func (s *Session) subPump() {
	defer func() {
		s.wg.Done()
		s.logger.Info("shutting down - stopping subPump")
	}()

	for {
		select {
		case c := <-s.checkChannel:
			request, ok := c.(*types.CheckRequest)
			if !ok {
				s.logger.Errorf("session received non-config over check channel")
				continue
			}

			configBytes, err := json.Marshal(request)
			if err != nil {
				s.logger.WithError(err).Error("session failed to serialize check request")
			}

			msg := &transport.Message{
				Type:    types.CheckRequestType,
				Payload: configBytes,
			}
			s.sendq <- msg
		case <-s.stopping:
			return
		}
	}
}

func (s *Session) sendPump() {
	defer func() {
		s.wg.Done()
		s.logger.Info("shutting down - stopping sendPump")
	}()

	for {
		select {
		case msg := <-s.sendq:
			s.logger.Debugf("session - sending message: %s", string(msg.Payload))
			err := s.conn.Send(msg)
			if err != nil {
				switch err := err.(type) {
				case transport.ConnectionError, transport.ClosedError:
					return
				default:
					s.logger.WithError(err).Error("send error")
				}
			}
		case <-s.stopping:
			return
		}
	}
}

// Start a Session.
// 1. Start send pump
// 2. Start receive pump
// 3. Start subscription pump
// 5. Ensure bus unsubscribe when the session shuts down.
func (s *Session) Start() error {
	s.wg = &sync.WaitGroup{}
	s.wg.Add(3)
	go s.sendPump()
	go s.recvPump()
	go s.subPump()

	for _, sub := range s.cfg.Subscriptions {
		topic := messaging.SubscriptionTopic(s.cfg.Organization, s.cfg.Environment, sub)
		s.logger.WithField("topic", topic).Debugf("subscribing to topic")
		if err := s.bus.Subscribe(topic, s.cfg.AgentID, s.checkChannel); err != nil {
			s.logger.WithError(err).Error("error starting subscription")
			return err
		}
	}

	return nil
}

// Stop a running session. This will cause the send and receive loops to
// shutdown. Blocks until the session has shutdown.
func (s *Session) Stop() {
	close(s.stopping)
	s.wg.Wait()
	for _, sub := range s.cfg.Subscriptions {
		s.bus.Unsubscribe(sub, s.ID)
	}
	close(s.checkChannel)
}

func (s *Session) handleKeepalive(payload []byte) error {
	keepalive := &types.Event{}
	err := json.Unmarshal(payload, keepalive)
	if err != nil {
		return err
	}

	// TODO(greg): better entity validation than this garbage.
	if keepalive.Entity == nil {
		return errors.New("keepalive does not contain an entity")
	}

	if keepalive.Timestamp == 0 {
		return errors.New("keepalive contains invalid timestamp")
	}

	return s.bus.Publish(messaging.TopicKeepalive, keepalive)
}

func (s *Session) handleEvent(payload []byte) error {
	event := &types.Event{}
	if err := json.Unmarshal(payload, event); err != nil {
		return err
	}
	return s.bus.Publish(messaging.TopicEventRaw, event)
}
