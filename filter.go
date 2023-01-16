package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// FilterFn represents the FilterIterator function.
type FilterFn[T any] func(v T) bool

// Filter creates a new FilterIterator.
func Filter[T, MAP any](iterator Iterator[T], fn FilterFn[T]) *FilterIterator[T, struct{}] {
	return FilterMap[T, struct{}](iterator, fn)
}

// FilterMap creates a new `FilterIterator` for use which also specifies a potential future `Map` operation.
func FilterMap[T, MAP any](iterator Iterator[T], fn FilterFn[T]) *FilterIterator[T, MAP] {
	return &FilterIterator[T, MAP]{
		iterator: iterator,
		fn:       fn,
	}
}

// FilterIterator allows filtering of an `Iterator[T]`.
type FilterIterator[T, MAP any] struct {
	iterator Iterator[T]
	fn       FilterFn[T]
}

// Next yields the next value from the iterator that passed the filter function.
func (i *FilterIterator[T, MAP]) Next() optionext.Option[T] {
	for {
		v := i.iterator.Next()
		if v.IsNone() || !i.fn(v.Unwrap()) {
			return v
		}
	}
}

// Iter is a convenience function that converts the `FilterIterator` iterator into an `*Iterate[T]`.
func (i *FilterIterator[T, MAP]) Iter() *Iterate[T, MAP] {
	return IterMap[T, MAP](i)
}
