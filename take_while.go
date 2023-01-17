package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// TakeWhileFn represents the `takeWhileIterator[T]` function.
type TakeWhileFn[T any] func(v T) bool

// TakeWhile creates a new `takeWhileIterator[T,I]` for use.
func TakeWhile[T any, I Iterator[T]](iterator I, fn TakeWhileFn[T]) takeWhileIterator[T, I, struct{}] {
	return TakeWhileWithMap[T, I, struct{}](iterator, fn)
}

// TakeWhileWithMap creates a new `takeWhileIterator[T,I]` for use and can specify a future `Map` type conversion.
func TakeWhileWithMap[T any, I Iterator[T], MAP any](iterator Iterator[T], fn TakeWhileFn[T]) takeWhileIterator[T, I, MAP] {
	return takeWhileIterator[T, I, MAP]{
		iterator: iterator,
		fn:       fn,
	}
}

// takeWhileIterator is an iterator that iterates over elements until the function return false or
// end of the iterator (whichever happens first).
type takeWhileIterator[T any, I Iterator[T], MAP any] struct {
	iterator Iterator[T]
	fn       TakeWhileFn[T]
}

// Next returns the next element until `TakeWhileFn[T]` returns false or end of the iterator.
func (i takeWhileIterator[T, I, MAP]) Next() optionext.Option[T] {
	for {
		v := i.iterator.Next()
		if v.IsNone() || i.fn(v.Unwrap()) {
			return v
		}
	}
}

// Iter is a convenience function that converts the `takeWhileIterator` iterator into an `*Iterate[T]`.
func (i takeWhileIterator[T, I, MAP]) Iter() Iterate[T, Iterator[T], MAP] {
	return IterMap[T, Iterator[T], MAP](i)
}
