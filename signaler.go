package syncache

import (
	"context"
)

// Operation is an Entry operation.
type Operation uint8

const (
	// OperationSet specifies that new entry has been set.
	OperationSet Operation = iota
	// OperationInvalidate specifies that entry has been invalidated.
	OperationInvalidate
)

type Signal[K comparable, V any] struct {
	Entry[K, V]
	Operation Operation
}

type Signaler[K comparable, V any] interface {
	// Signal signals that entries where mutated in external datasource and needs to be mutated in caches interested
	// in maintain the same eventually-consistent state too.
	Signal(ctx context.Context, signals ...Signal[K, V]) error
}

func NewSignal[K comparable, V any](entry Entry[K, V], operation Operation) Signal[K, V] {
	return Signal[K, V]{
		Entry:     entry,
		Operation: operation,
	}
}
