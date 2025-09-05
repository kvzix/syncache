package syncache

import "time"

type (
	InitialGetOptions[K comparable, V any] struct {
		Loader      Loader[K, V]
		BatchLoader BatchLoader[K, V]
	}

	GetOptionsStore[K comparable, V any] struct {
		loader      Loader[K, V]
		batchLoader BatchLoader[K, V]
	}
)

type GetOption[K comparable, V any] func(*GetOptionsStore[K, V])

// buildGetOptions builds and initialize GetOptionsStore based on initial and external options.
func buildGetOptions[K comparable, V any](initialOptions InitialGetOptions[K, V], options []GetOption[K, V]) GetOptionsStore[K, V] {
	store := GetOptionsStore[K, V]{
		loader:      initialOptions.Loader,
		batchLoader: initialOptions.BatchLoader,
	}

	for _, apply := range options {
		apply(&store)
	}

	return store
}

// WithLoader sets one-time Loader that will be used instead of Cache's internal loader (if any).
func WithLoader[K comparable, V any](loader Loader[K, V]) GetOption[K, V] {
	return func(o *GetOptionsStore[K, V]) {
		o.loader = loader
	}
}

// Loader returns Loader set by WithLoader.
func (os GetOptionsStore[K, V]) Loader() Loader[K, V] {
	return os.loader
}

// WithBatchLoader sets one-time BatchLoader that will be used instead of Cache's internal loader (if any).
func WithBatchLoader[K comparable, V any](batchLoader BatchLoader[K, V]) GetOption[K, V] {
	return func(o *GetOptionsStore[K, V]) {
		o.batchLoader = batchLoader
	}
}

// BatchLoader returns BatchLoader set by WithBatchLoader.
func (os GetOptionsStore[K, V]) BatchLoader() BatchLoader[K, V] {
	return os.batchLoader
}

type (
	InitialSetOptions[K comparable, V any] struct {
		TTL time.Duration
	}

	SetOptionsStore[K comparable, V any] struct {
		ttl time.Duration
	}
)

type SetOption[K comparable, V any] func(*SetOptionsStore[K, V])

// buildSetOptions builds and initialize SetOptionsStore based on initial and external options.
func buildSetOptions[K comparable, V any](initialOptions InitialSetOptions[K, V], options []SetOption[K, V]) SetOptionsStore[K, V] {
	store := SetOptionsStore[K, V]{
		ttl: initialOptions.TTL,
	}

	for _, apply := range options {
		apply(&store)
	}

	return store
}

func WithTTL[K comparable, V any](ttl time.Duration) SetOption[K, V] {
	return func(o *SetOptionsStore[K, V]) {
		o.ttl = ttl
	}
}

func (os SetOptionsStore[K, V]) TTL() time.Duration {
	return os.ttl
}
