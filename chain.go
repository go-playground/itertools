package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Chain creates a new `ChainIterator[T]` for use.
func Chain[T any](first, second Iterator[T]) *ChainIterator[T] {
	return &ChainIterator[T]{
		current: first,
		next:    second,
	}
}

// ChainIterator takes two iterators and creates a new iterator over both in sequence.
type ChainIterator[T any] struct {
	current Iterator[T]
	next    Iterator[T]
	flipped bool
}

// Next returns the next value from the first iterator until exhausted and then the second.
func (i *ChainIterator[T]) Next() optionext.Option[T] {
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
