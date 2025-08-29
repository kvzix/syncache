package syncache

import (
	"context"
	"fmt"
)

// Mutator is a function that mutates entry in the external datasource and returns it in mutated state.
type Mutator[K comparable, V any] func() (Signal[K, V], error)

// Mutate mutates entry with Mutator and signalize to interested caches that they should mutate entry returned from Mutator.
func Mutate[K comparable, V any](ctx context.Context, signaler Signaler[K, V], mutate Mutator[K, V]) error {
	return MutateBatch(ctx, signaler, func() ([]Signal[K, V], error) {
		mutatedEntry, err := mutate()
		if err != nil {
			return nil, err
		}

		return []Signal[K, V]{mutatedEntry}, nil
	})
}

// BatchMutator is a function that mutates entries in the external datasource and returns them in mutated state.
type BatchMutator[K comparable, V any] func() ([]Signal[K, V], error)

// MutateBatch mutates entries BatchMutator and signalize to interested caches that they should mutate entries returned from BatchMutator.
func MutateBatch[K comparable, V any](ctx context.Context, signaler Signaler[K, V], mutate BatchMutator[K, V]) error {
	mutatedEntries, err := mutate()
	if err != nil {
		return err
	}

	// Signal that entries were mutated. Usually, if you have multiple instances of your application it can be done via some
	// queue like Nats, RabbitMQ etc.
	if err = signaler.Signal(ctx, mutatedEntries...); err != nil {
		return fmt.Errorf("signal mutations: %w", err)
	}

	return nil
}
