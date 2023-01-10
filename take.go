package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Take creates a new `TakeIterator[T]` for use.
func Take[T any](iterator Iterator[T], n int) *TakeIterator[T] {
	return &TakeIterator[T]{
		iterator: iterator,
		limit:    n,
	}
}

// TakeIterator is an iterator that only iterates over n elements.
type TakeIterator[T any] struct {
	iterator Iterator[T]
	limit    int
}

// Next returns the next element until n is reached or end of the iterator.
func (i *TakeIterator[T]) Next() optionext.Option[T] {
	if i.limit <= 0 {
		return optionext.None[T]()
	}
	i.limit--
	return i.iterator.Next()
}
