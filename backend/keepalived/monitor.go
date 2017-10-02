package keepalived

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sensu/sensu-go/backend/store"
	"github.com/sensu/sensu-go/types"
)

// KeepaliveMonitor is a managed timer that is reset whenever the monitor
// observes a Keepalive event via the Update() function.
type KeepaliveMonitor struct {
	Entity       *types.Entity
	Deregisterer Deregisterer
	EventCreator EventCreator
	Store        store.Store

	reset   chan interface{}
	timer   *time.Timer
	stopped int32
	failing int32

	mu sync.Mutex
}

// Start initializes the monitor and starts its monitoring goroutine.
func (monitorPtr *KeepaliveMonitor) Start() {
	timerDuration := time.Duration(monitorPtr.Entity.KeepaliveTimeout) * time.Second
	monitorPtr.timer = time.NewTimer(timerDuration)
	monitorPtr.reset = make(chan interface{})
	go func() {
		timer := monitorPtr.timer
		ctx := context.WithValue(context.Background(), types.OrganizationKey, monitorPtr.Entity.Organization)
		ctx = context.WithValue(ctx, types.EnvironmentKey, monitorPtr.Entity.Environment)

		var (
			event   *types.Event
			err     error
			timeout int64
		)

		for {
			select {
			case <-monitorPtr.reset:
				monitorPtr.mu.Lock()

				if !timer.Stop() {
					<-timer.C
				}
				if monitorPtr.IsStopped() {
					monitorPtr.mu.Unlock()
					return
				}
				monitorPtr.mu.Unlock()

			case <-timer.C:
				// timed out keepalive

				// test to see if the entity still exists (it may have been deleted)

				event, err = monitorPtr.Store.GetEventByEntityCheck(ctx, monitorPtr.Entity.ID, "keepalive")
				if err != nil {
					// this should be a temporary error talking to the store. keep trying until
					// the store starts responding again.
					logger.WithError(err).Error("error getting keepalive event for client")
					break
				}

				// if the agent disconnected and reconnected elsewhere, stop the monitor
				// and return.
				if event != nil && event.Check.Status == 0 {
					monitorPtr.Store.DeleteFailingKeepalive(ctx, monitorPtr.Entity)
					monitorPtr.Stop()
					return
				}

				// if the entity is supposed to be deregistered, do so.
				if monitorPtr.Entity.Deregister {
					if err = monitorPtr.Deregisterer.Deregister(monitorPtr.Entity); err != nil {
						logger.WithError(err).Error("error deregistering entity")
					}
					monitorPtr.Stop()
					return
				}

				// this is a real keepalive event, emit it.
				if err = monitorPtr.EventCreator.Warn(monitorPtr.Entity); err != nil {
					logger.WithError(err).Error("error sending keepalive event")
				}

				timeout = time.Now().Unix() + int64(monitorPtr.Entity.KeepaliveTimeout)
				if err = monitorPtr.Store.UpdateFailingKeepalive(ctx, monitorPtr.Entity, timeout); err != nil {
					logger.WithError(err).Error("error updating failing keepalive in store")
				}

				atomic.CompareAndSwapInt32(&monitorPtr.failing, 0, 1)
			}

			monitorPtr.mu.Lock()
			timer.Reset(timerDuration)
			monitorPtr.mu.Unlock()
		}
	}()
}

// Update causes the KeepaliveMonitor to observe the event.
func (monitorPtr *KeepaliveMonitor) Update(event *types.Event) error {
	entity := event.Entity

	if atomic.CompareAndSwapInt32(&monitorPtr.failing, 1, 0) {
		monitorPtr.Store.DeleteFailingKeepalive(context.Background(), entity)
	}

	monitorPtr.reset <- struct{}{}

	entity.LastSeen = event.Timestamp
	ctx := context.WithValue(context.Background(), types.OrganizationKey, entity.Organization)
	ctx = context.WithValue(ctx, types.EnvironmentKey, monitorPtr.Entity.Environment)

	if err := monitorPtr.Store.UpdateEntity(ctx, entity); err != nil {
		logger.WithError(err).Error("error updating entity in store")
	}

	return monitorPtr.EventCreator.Pass(entity)
}

// Stop the KeepaliveMonitor
func (monitorPtr *KeepaliveMonitor) Stop() {
	// atomically set stopped so that once Stop is called, all future
	// reads of stopped are true.
	if !atomic.CompareAndSwapInt32(&monitorPtr.stopped, 0, 1) {
		return
	}

	close(monitorPtr.reset)
}

// IsStopped returns true if the Monitor has been stopped.
func (monitorPtr *KeepaliveMonitor) IsStopped() bool {
	return atomic.LoadInt32(&monitorPtr.stopped) > 0
}

// Reset the monitor's timer to emit an event at a given time.
func (monitorPtr *KeepaliveMonitor) Reset(t int64) {
	monitorPtr.mu.Lock()
	defer monitorPtr.mu.Unlock()

	if monitorPtr.timer == nil {
		monitorPtr.Start()
	}

	d := time.Duration(t - time.Now().Unix())
	if d < 0 {
		d = 0
	}

	monitorPtr.timer.Reset(d)
}
