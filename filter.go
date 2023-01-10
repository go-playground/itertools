package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// FilterFn represents the FilterIterator function.
type FilterFn[T any] func(v T) bool

// Filter creates a new FilterIterator.
func Filter[T any](iterator Iterator[T], fn FilterFn[T]) *FilterIterator[T] {
	return &FilterIterator[T]{
		iterator: iterator,
		fn:       fn,
	}
}

// FilterIterator allows filtering of an `Iterator[T]`.
type FilterIterator[T any] struct {
	iterator Iterator[T]
	fn       FilterFn[T]
}

// Next yields the next value from the iterator that passed the filter function.
func (i *FilterIterator[T]) Next() optionext.Option[T] {
	for {
		v := i.iterator.Next()
		if v.IsNone() || !i.fn(v.Unwrap()) {
			return v
		}
	}
}
