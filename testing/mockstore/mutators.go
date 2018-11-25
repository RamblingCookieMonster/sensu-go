package mockstore

import (
	"context"

	types "github.com/sensu/sensu-go/api/core/v2"
)

// DeleteMutatorByName ...
func (s *MockStore) DeleteMutatorByName(ctx context.Context, name string) error {
	args := s.Called(ctx, name)
	return args.Error(0)
}

// GetMutators ...
func (s *MockStore) GetMutators(ctx context.Context) ([]*types.Mutator, error) {
	args := s.Called(ctx)
	return args.Get(0).([]*types.Mutator), args.Error(1)
}

// GetMutatorByName ...
func (s *MockStore) GetMutatorByName(ctx context.Context, name string) (*types.Mutator, error) {
	args := s.Called(ctx, name)
	return args.Get(0).(*types.Mutator), args.Error(1)
}

// UpdateMutator ...
func (s *MockStore) UpdateMutator(ctx context.Context, mutator *types.Mutator) error {
	args := s.Called(mutator)
	return args.Error(0)
}
