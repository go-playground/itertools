package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// FilterFn represents the `filterIterator` function.
type FilterFn[T any] func(v T) bool

// Filter creates a new `filterIterator`.
func Filter[T, MAP any](iterator Iterator[T], fn FilterFn[T]) *filterIterator[T, struct{}] {
	return FilterWithMap[T, struct{}](iterator, fn)
}

// FilterWithMap creates a new `filterIterator` for use which also specifies a potential future `Map` operation.
func FilterWithMap[T, MAP any](iterator Iterator[T], fn FilterFn[T]) *filterIterator[T, MAP] {
	return &filterIterator[T, MAP]{
		iterator: iterator,
		fn:       fn,
	}
}

// filterIterator allows filtering of an `Iterator[T]`.
type filterIterator[T, MAP any] struct {
	iterator Iterator[T]
	fn       FilterFn[T]
}

// Next yields the next value from the iterator that passed the filter function.
func (i *filterIterator[T, MAP]) Next() optionext.Option[T] {
	for {
		v := i.iterator.Next()
		if v.IsNone() || !i.fn(v.Unwrap()) {
			return v
		}
	}
}

// Iter is a convenience function that converts the `filterIterator` iterator into an `*Iterate[T]`.
func (i *filterIterator[T, MAP]) Iter() Iterate[T, MAP] {
	return IterMap[T, MAP](i)
}
