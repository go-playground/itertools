package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// TakeWhileFn represents the `takeWhileIterator[T]` function.
type TakeWhileFn[T any] func(v T) bool

// TakeWhile creates a new `takeWhileIterator[T,I]` for use.
func TakeWhile[T any](iterator Iterator[T], fn TakeWhileFn[T]) takeWhileIterator[T, struct{}] {
	return TakeWhileMap[T, struct{}](iterator, fn)
}

// TakeWhileMap creates a new `takeWhileIterator[T,I]` for use and can specify a future `Map` type conversion.
func TakeWhileMap[T, MAP any](iterator Iterator[T], fn TakeWhileFn[T]) takeWhileIterator[T, MAP] {
	return takeWhileIterator[T, MAP]{
		iterator: iterator,
		fn:       fn,
	}
}

// takeWhileIterator is an iterator that iterates over elements until the function return false or
// end of the iterator (whichever happens first).
type takeWhileIterator[T, MAP any] struct {
	iterator Iterator[T]
	fn       TakeWhileFn[T]
}

// Next returns the next element until `TakeWhileFn[T]` returns false or end of the iterator.
func (i takeWhileIterator[T, MAP]) Next() optionext.Option[T] {
	for {
		v := i.iterator.Next()
		if v.IsNone() || i.fn(v.Unwrap()) {
			return v
		}
	}
}

// Iter is a convenience function that converts the `takeWhileIterator` iterator into an `*Iterate[T]`.
func (i takeWhileIterator[T, MAP]) Iter() Iterate[T, MAP] {
	return IterMap[T, MAP](i)
}
