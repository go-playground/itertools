package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Peekable accepts and `Iterator[T]` and turns it into a Peekable iterator.
//
// NOTE: Peekable iterators are commonly the LAST in a chain of iterators.
func Peekable[T any, I Iterator[T]](iterator I) *peekableIterator[T, I] {
	return &peekableIterator[T, I]{
		iterator: iterator,
	}
}

// peekableIterator makes an `Iterator` peekable.
type peekableIterator[T any, I Iterator[T]] struct {
	iterator I
	prev     optionext.Option[T]
}

// Next advances the iterator and returns the next value.
//
// Returns an Option with value of None when iteration has finished.
func (i *peekableIterator[T, I]) Next() optionext.Option[T] {
	if i.prev.IsSome() {
		prev := i.prev
		i.prev = optionext.None[T]()
		return prev
	}
	return i.iterator.Next()
}

// Peek returns the next value without advancing the iterator.
func (i *peekableIterator[T, I]) Peek() optionext.Option[T] {
	if i.prev.IsSome() {
		return i.prev
	}
	i.prev = i.iterator.Next()
	return i.prev
}
