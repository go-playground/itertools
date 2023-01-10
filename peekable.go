package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Peekable accepts and `Iterator[T]` and turns it into a Peekable iterator.
//
// NOTE: Peekable iterators are commonly the LAST in a chain of iterators.
func Peekable[T any](iterator Iterator[T]) *PeekableIterator[T] {
	return &PeekableIterator[T]{
		iterator: iterator,
	}
}

// PeekableIterator makes an `Iterator` peekable.
type PeekableIterator[T any] struct {
	iterator Iterator[T]
	prev     optionext.Option[T]
}

// Next advances the iterator and returns the next value.
//
// Returns an Option with value of None when iteration has finished.
func (i *PeekableIterator[T]) Next() optionext.Option[T] {
	if i.prev.IsSome() {
		prev := i.prev
		i.prev = optionext.None[T]()
		return prev
	}
	return i.iterator.Next()
}

// Peek returns the next value without advancing the iterator.
func (i *PeekableIterator[T]) Peek() optionext.Option[T] {
	if i.prev.IsSome() {
		return i.prev
	}
	i.prev = i.iterator.Next()
	return i.prev
}
