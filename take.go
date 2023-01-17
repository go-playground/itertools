package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Take creates a new `takeIterator[T]` for use.
func Take[T any, I Iterator[T]](iterator I, n int) *takeIterator[T, I, struct{}] {
	return TakeWithMap[T, I, struct{}](iterator, n)
}

// TakeWithMap creates a new `takeIterator[T]` for use and can specify a future `Map` type conversion.
func TakeWithMap[T any, I Iterator[T], MAP any](iterator I, n int) *takeIterator[T, I, MAP] {
	return &takeIterator[T, I, MAP]{
		iterator: iterator,
		limit:    n,
	}
}

// takeIterator is an iterator that only iterates over n elements.
type takeIterator[T any, I Iterator[T], MAP any] struct {
	iterator I
	limit    int
}

// Next returns the next element until n is reached or end of the iterator.
func (i *takeIterator[T, I, MAP]) Next() optionext.Option[T] {
	if i.limit <= 0 {
		return optionext.None[T]()
	}
	i.limit--
	return i.iterator.Next()
}

// Iter is a convenience function that converts the `takeIterator` iterator into an `*Iterate[T]`.
func (i *takeIterator[T, I, MAP]) Iter() Iterate[T, Iterator[T], MAP] {
	return IterMap[T, Iterator[T], MAP](i)
}
