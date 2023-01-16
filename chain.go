package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Chain creates a new `chainIterator[T]` for use.
func Chain[T any](first, second Iterator[T]) *chainIterator[T, struct{}] {
	return ChainMap[T, struct{}](first, second)
}

// ChainMap creates a new `chainIterator[T]` for use and parameter to specify a Map type for the `Iterate.Map` helper
// function.
func ChainMap[T, MAP any](first, second Iterator[T]) *chainIterator[T, MAP] {
	return &chainIterator[T, MAP]{
		current: first,
		next:    second,
	}
}

// chainIterator takes two iterators and creates a new iterator over both in sequence.
type chainIterator[T, MAP any] struct {
	current Iterator[T]
	next    Iterator[T]
	flipped bool
}

// Next returns the next value from the first iterator until exhausted and then the second.
func (i *chainIterator[T, MAP]) Next() optionext.Option[T] {
	for {
		v := i.current.Next()
		if v.IsSome() {
			return v
		}
		if i.flipped {
			return v
		}
		i.current = i.next
		i.flipped = true
	}
}

// Iter is a convenience function that converts the chainIterator iterator into an `*Iterate[T]`.
func (i *chainIterator[T, MAP]) Iter() *Iterate[T, MAP] {
	return IterMap[T, MAP](i)
}
