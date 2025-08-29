package syncache

type receiverOptionStore[K comparable, V any] struct {
	applier Applier[K, V]
}

type ReceiverOption[K comparable, V any] func(store *receiverOptionStore[K, V])

func WithApplier[K comparable, V any](applier Applier[K, V]) ReceiverOption[K, V] {
	return func(r *receiverOptionStore[K, V]) {
		r.applier = applier
	}
}
