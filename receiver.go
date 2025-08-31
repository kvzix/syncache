package syncache

import (
	"context"
	"errors"
	"fmt"
)

// ErrEmptyReceiver indicates that channel returned from Receiver is empty. This error is used to prevent deadlock.
var ErrEmptyReceiver = errors.New("receiver channel is empty")

// Applier applies signals from Receiver to Cache.
type Applier[K comparable, V any] func(ctx context.Context, cache Cache[K, V], signals ...Signal[K, V]) error

type Receiver[K comparable, V any] interface {
	Receive(ctx context.Context) (<-chan []Signal[K, V], error)
}

func Run[K comparable, V any](ctx context.Context, cache Cache[K, V], receiver Receiver[K, V], options ...ReceiverOption[K, V]) error {
	var optionStore receiverOptionStore[K, V]

	for _, option := range options {
		option(&optionStore)
	}

	signals, err := receiver.Receive(ctx)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}

	if signals == nil {
		return ErrEmptyReceiver
	}

	applier := optionStore.applier
	if applier == nil {
		applier = applySignals
	}

	for entries := range signals {
		if err = applier(ctx, cache, entries...); err != nil {
			return fmt.Errorf("apply signals: %w", err)
		}
	}

	return nil
}

func applySignals[K comparable, V any](ctx context.Context, cache Cache[K, V], signals ...Signal[K, V]) error {
	var (
		sets        []Entry[K, V]
		invalidates []K
	)

	for _, signalEntry := range signals {
		operation := signalEntry.Operation

		switch operation {
		case OperationSet:
			sets = append(sets, signalEntry.Entry)
		case OperationInvalidate:
			invalidates = append(invalidates, signalEntry.Entry.Key)
		}
	}

	if _, err := setByLength(ctx, cache, sets...); err != nil {
		return err
	}

	if _, err := invalidateByLength(ctx, cache, invalidates...); err != nil {
		return err
	}

	return nil
}
