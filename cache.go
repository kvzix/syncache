package syncache

import (
	"context"
	"fmt"
)

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

type (
	Loader[K comparable, V any]      func(ctx context.Context, key K) (V, error)
	BatchLoader[K comparable, V any] func(ctx context.Context) ([]Entry[K, V], error)
)

type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K, options ...CacheOption[K, V]) (V, error)
	GetBatch(ctx context.Context, keys []K, options ...CacheOption[K, V]) ([]Entry[K, V], error)
	Set(ctx context.Context, entry Entry[K, V]) error
	SetBatch(ctx context.Context, entries []Entry[K, V]) error
	Invalidate(ctx context.Context, key K) error
	InvalidateBatch(ctx context.Context, keys []K) error
}

func NewEntry[K comparable, V any](key K, value V) Entry[K, V] {
	return Entry[K, V]{
		Key:   key,
		Value: value,
	}
}

// Load loads values in batch in Cache with entries from BatchLoader.
//
// This helper may be useful when you want to warm your Cache on start of application with entries from the external datasource.
func Load[K comparable, V any](ctx context.Context, cache Cache[K, V], loadBatch BatchLoader[K, V]) error {
	values, err := loadBatch(ctx)
	if err != nil {
		return fmt.Errorf("load batch: %w", err)
	}

	if err = cache.SetBatch(ctx, values); err != nil {
		return fmt.Errorf("set batch: %w", err)
	}

	return nil
}
