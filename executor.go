package syncache

import (
	"context"
	"errors"
)

// ErrNoSource indicates that cache is disabled via WithoutCache but no Loader and/or BatchLoader is provided so there is
// no source to take data from.
var ErrNoSource = errors.New("data source does not provided")

type (
	Executor[K comparable, V any] struct {
		getOptions InitialGetOptions[K, V]
		setOptions InitialSetOptions[K, V]
	}

	ExecutorOptions[K comparable, V any] struct {
		GetOptions InitialGetOptions[K, V]
		SetOptions InitialSetOptions[K, V]
	}
)

func NewExecutor[K comparable, V any](options ExecutorOptions[K, V]) Executor[K, V] {
	return Executor[K, V]{
		getOptions: options.GetOptions,
		setOptions: options.SetOptions,
	}
}

func (e Executor[K, V]) Get(
	ctx context.Context,
	key K,
	options []GetOption[K, V],
	get func(GetOptionsStore[K, V]) (V, error),
) (V, error) {
	store := buildGetOptions(e.getOptions, options)

	if !isCacheGetterDisabled(ctx) {
		return get(store)
	}

	var value V

	load := store.Loader()
	if load == nil {
		return value, ErrNoSource
	}

	return load(ctx, key)
}

func (e Executor[K, V]) GetBatch(
	ctx context.Context,
	keys []K,
	options []GetOption[K, V],
	getBatch func(GetOptionsStore[K, V]) ([]Entry[K, V], error),
) ([]Entry[K, V], error) {
	store := buildGetOptions(e.getOptions, options)

	if !isCacheGetterDisabled(ctx) {
		return getBatch(store)
	}

	loadBatch := store.BatchLoader()
	if loadBatch == nil {
		return nil, ErrNoSource
	}

	return loadBatch(ctx, keys...)
}

func (e Executor[K, V]) Set(
	ctx context.Context,
	options []SetOption[K, V],
	set func(SetOptionsStore[K, V]) error,
) error {
	if isCacheSetterDisabled(ctx) {
		return nil
	}

	store := buildSetOptions(e.setOptions, options)
	return set(store)
}

func (e Executor[K, V]) Invalidate(ctx context.Context, invalidate func() error) error {
	if isCacheInvalidatorDisabled(ctx) {
		return nil
	}

	return invalidate()
}
