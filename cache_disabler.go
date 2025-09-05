package syncache

import (
	"context"
)

type DisableCacheOption func(*disableCacheValue)

type (
	disableCacheKey   struct{}
	disableCacheValue struct {
		get        bool
		set        bool
		invalidate bool
	}
)

func WithCacheDisabler(ctx context.Context, options ...DisableCacheOption) context.Context {
	var value disableCacheValue

	for _, option := range options {
		option(&value)
	}

	// Disable all operations if user don't specified specific ones.
	if !value.get && !value.set && !value.invalidate {
		value = disableCacheValue{
			get:        true,
			set:        true,
			invalidate: true,
		}
	}

	return context.WithValue(ctx, disableCacheKey{}, value)
}

func DisableGetter() DisableCacheOption {
	return func(value *disableCacheValue) {
		value.get = true
	}
}

func isCacheGetterDisabled(ctx context.Context) bool {
	if value, ok := ctx.Value(disableCacheKey{}).(disableCacheValue); ok {
		return value.get
	}

	return false
}

func DisableSetter() DisableCacheOption {
	return func(value *disableCacheValue) {
		value.set = true
	}
}

func isCacheSetterDisabled(ctx context.Context) bool {
	if value, ok := ctx.Value(disableCacheKey{}).(disableCacheValue); ok {
		return value.set
	}

	return false
}

func DisableInvalidator() DisableCacheOption {
	return func(value *disableCacheValue) {
		value.invalidate = true
	}
}

func isCacheInvalidatorDisabled(ctx context.Context) bool {
	if value, ok := ctx.Value(disableCacheKey{}).(disableCacheValue); ok {
		return value.invalidate
	}

	return false
}
