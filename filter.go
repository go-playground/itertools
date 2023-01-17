package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// FilterFn represents the `filterIterator` function.
type FilterFn[T any] func(v T) bool

// Filter creates a new `filterIterator`.
func Filter[T any, I Iterator[T]](iterator I, fn FilterFn[T]) *filterIterator[T, I, struct{}] {
	return FilterWithMap[T, I, struct{}](iterator, fn)
}

// FilterWithMap creates a new `filterIterator` for use which also specifies a potential future `Map` operation.
func FilterWithMap[T any, I Iterator[T], MAP any](iterator I, fn FilterFn[T]) *filterIterator[T, I, MAP] {
	return &filterIterator[T, I, MAP]{
		iterator: iterator,
		fn:       fn,
	}
}

// filterIterator allows filtering of an `Iterator[T]`.
type filterIterator[T any, I Iterator[T], MAP any] struct {
	iterator I
	fn       FilterFn[T]
}

// Next yields the next value from the iterator that passed the filter function.
func (i *filterIterator[T, I, MAP]) Next() optionext.Option[T] {
	for {
		v := i.iterator.Next()
		if v.IsNone() || !i.fn(v.Unwrap()) {
			return v
		}
	}
}

// Iter is a convenience function that converts the `filterIterator` iterator into an `Iterate[T]`.
func (i *filterIterator[T, I, MAP]) Iter() Iterate[T, Iterator[T], MAP] {
	return IterMap[T, Iterator[T], MAP](i)
}
