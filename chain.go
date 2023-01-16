package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Chain creates a new `chainIterator[T]` for use.
func Chain[T any, FI Iterator[T], SI Iterator[T]](first FI, second SI) *chainIterator[T, FI, SI, struct{}] {
	return ChainWithMap[T, FI, SI, struct{}](first, second)
}

// ChainWithMap creates a new `chainIterator[T]` for use and parameter to specify a Map type for the `Iterate.Map` helper
// function.
func ChainWithMap[T any, FI Iterator[T], SI Iterator[T], MAP any](first FI, second SI) *chainIterator[T, FI, SI, MAP] {
	return &chainIterator[T, FI, SI, MAP]{
		current: first,
		next:    second,
	}
}

// chainIterator takes two iterators and creates a new iterator over both in sequence.
type chainIterator[T any, FI Iterator[T], SI Iterator[T], MAP any] struct {
	current FI
	next    SI
	flipped bool
}

// Next returns the next value from the first iterator until exhausted and then the second.
func (i *chainIterator[T, FI, SI, MAP]) Next() optionext.Option[T] {
	for {
		if i.flipped {
			return i.next.Next()
		}
		v := i.current.Next()
		if v.IsSome() {
			return v
		}
		i.flipped = true
	}
}

// Iter is a convenience function that converts the chainIterator iterator into an `*Iterate[T]`.
func (i *chainIterator[T, FI, SI, MAP]) Iter() Iterate[T, Iterator[T], MAP] {
	return IterMap[T, Iterator[T], MAP](i)
}
