package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Take creates a new `takeIterator[T]` for use.
func Take[T any](iterator Iterator[T], n int) *takeIterator[T, struct{}] {
	return TakeMap[T, struct{}](iterator, n)
}

// TakeMap creates a new `takeIterator[T]` for use and can specify a future `Map` type conversion.
func TakeMap[T, MAP any](iterator Iterator[T], n int) *takeIterator[T, MAP] {
	return &takeIterator[T, MAP]{
		iterator: iterator,
		limit:    n,
	}
}

// takeIterator is an iterator that only iterates over n elements.
type takeIterator[T, MAP any] struct {
	iterator Iterator[T]
	limit    int
}

// Next returns the next element until n is reached or end of the iterator.
func (i *takeIterator[T, MAP]) Next() optionext.Option[T] {
	if i.limit <= 0 {
		return optionext.None[T]()
	}
	i.limit--
	return i.iterator.Next()
}

// Iter is a convenience function that converts the `takeIterator` iterator into an `*Iterate[T]`.
func (i *takeIterator[T, MAP]) Iter() Iterate[T, MAP] {
	return IterMap[T, MAP](i)
}
