package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// TakeWhileFn represents the `TakeWhileIterator[T]` function.
type TakeWhileFn[T any] func(v T) bool

// TakeWhile creates a new `TakeWhileIterator[T,I]` for use.
func TakeWhile[T any](iterator Iterator[T], fn TakeWhileFn[T]) *TakeWhileIterator[T] {
	return &TakeWhileIterator[T]{
		iterator: iterator,
		fn:       fn,
	}
}

// TakeWhileIterator is an iterator that iterates over elements until the function return false or
// end of the iterator (whichever happens first).
type TakeWhileIterator[T any] struct {
	iterator Iterator[T]
	fn       TakeWhileFn[T]
}

// Next returns the next element until `TakeWhileFn[T]` returns false or end of the iterator.
func (i *TakeWhileIterator[T]) Next() optionext.Option[T] {
	for {
		v := i.iterator.Next()
		if v.IsNone() || i.fn(v.Unwrap()) {
			return v
		}
	}
}
