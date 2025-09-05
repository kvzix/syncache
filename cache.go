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
	// Loader is a function that loads value by key from the external datasource.
	Loader[K comparable, V any] func(ctx context.Context, key K) (V, error)
	// BatchLoader ...
	BatchLoader[K comparable, V any] func(ctx context.Context, keys ...K) ([]Entry[K, V], error)
	// SnapshotLoader ...
	SnapshotLoader[K comparable, V any] func(ctx context.Context) ([]Entry[K, V], error)
)

type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K, options ...GetOption[K, V]) (V, error)
	GetBatch(ctx context.Context, keys []K, options ...GetOption[K, V]) ([]Entry[K, V], error)
	Set(ctx context.Context, entry Entry[K, V], options ...SetOption[K, V]) error
	SetBatch(ctx context.Context, entries []Entry[K, V], options ...SetOption[K, V]) error
	Invalidate(ctx context.Context, key K) error
	InvalidateBatch(ctx context.Context, keys []K) error
}

func NewEntry[K comparable, V any](key K, value V) Entry[K, V] {
	return Entry[K, V]{
		Key:   key,
		Value: value,
	}
}

// Load loads entries from SnapshotLoader in Cache.
//
// This helper function may be useful if you want to warm your Cache on start of application with entries from the external datasource.
func Load[K comparable, V any](ctx context.Context, cache Cache[K, V], loadSnapshot SnapshotLoader[K, V]) error {
	loadedValues, err := loadSnapshot(ctx)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	// Set values loaded by batch loader in cache.
	if _, err = setByLength(ctx, cache, loadedValues...); err != nil {
		return err
	}

	return nil
}

func setByLength[K comparable, V any](ctx context.Context, cache Cache[K, V], entries ...Entry[K, V]) (bool, error) {
	switch len(entries) {
	case 0:
		return false, nil
	case 1:
		if err := cache.Set(ctx, entries[0]); err != nil {
			return false, fmt.Errorf("set: %w", err)
		}
	}

	if err := cache.SetBatch(ctx, entries); err != nil {
		return false, fmt.Errorf("set batch: %w", err)
	}

	return true, nil
}

func invalidateByLength[K comparable, V any](ctx context.Context, cache Cache[K, V], keys ...K) (bool, error) {
	switch len(keys) {
	case 0:
		return false, nil
	case 1:
		if err := cache.Invalidate(ctx, keys[0]); err != nil {
			return false, fmt.Errorf("invalidate: %w", err)
		}
	}

	if err := cache.InvalidateBatch(ctx, keys); err != nil {
		return false, fmt.Errorf("invalidate batch: %w", err)
	}

	return true, nil
}
