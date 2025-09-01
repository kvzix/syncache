package syncache

import (
	"context"
	"fmt"
)

type (
	// MutatorFunc is a function that mutates entry in the external datasource and returns it in mutated state.
	MutatorFunc[K comparable, V any] func() (Signal[K, V], error)
	// BatchMutatorFunc is a function that mutates entries in the external datasource and returns them in mutated state.
	BatchMutatorFunc[K comparable, V any] func() ([]Signal[K, V], error)
)

type Mutator[K comparable, V any] struct {
	signaler Signaler[K, V]
}

func NewMutator[K comparable, V any](signaler Signaler[K, V]) Mutator[K, V] {
	return Mutator[K, V]{
		signaler: signaler,
	}
}

// Mutate mutates entry with MutatorFunc and signalize to interested caches that they should mutate entry returned from Mutator.
func (m Mutator[K, V]) Mutate(ctx context.Context, mutate MutatorFunc[K, V]) error {
	return m.MutateBatch(ctx, func() ([]Signal[K, V], error) {
		mutatedEntry, err := mutate()
		if err != nil {
			return nil, err
		}

		return []Signal[K, V]{mutatedEntry}, nil
	})
}

// MutateBatch mutates entries BatchMutator and signalize to interested caches that they should mutate entries returned from BatchMutator.
func (m Mutator[K, V]) MutateBatch(ctx context.Context, mutate BatchMutatorFunc[K, V]) error {
	mutatedEntries, err := mutate()
	if err != nil {
		return err
	}

	// Signal that entries were mutated.
	if err = m.signaler.Signal(ctx, mutatedEntries...); err != nil {
		return fmt.Errorf("signal mutations: %w", err)
	}

	return nil
}
