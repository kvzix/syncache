package syncache

type CacheOptionStore[K comparable, V any] struct {
	loader Loader[K, V]
}

type CacheOption[K comparable, V any] func(*CacheOptionStore[K, V])

// Options returns CacheOptionStore composed from the provided options.
func Options[K comparable, V any](options []CacheOption[K, V]) CacheOptionStore[K, V] {
	var store CacheOptionStore[K, V]

	for _, apply := range options {
		apply(&store)
	}

	return store
}

// WithLoader sets one-time Loader that will be used instead of Cache's internal loader (if any).
func WithLoader[K comparable, V any](loader Loader[K, V]) CacheOption[K, V] {
	return func(o *CacheOptionStore[K, V]) {
		o.loader = loader
	}
}

// Loader returns Loader set by WithLoader.
func (os CacheOptionStore[K, V]) Loader() Loader[K, V] {
	return os.loader
}
