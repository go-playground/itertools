package itertools

import optionext "github.com/go-playground/pkg/v5/values/option"

// SliceIter accepts and turns a slice into an iterator
func SliceIter[T any](slice []T) *sliceIterator[T] {
	return &sliceIterator[T]{
		slice: slice,
	}
}

type sliceIterator[T any] struct {
	slice []T
}

func (i *sliceIterator[T]) Next() optionext.Option[T] {
	if len(i.slice) == 0 {
		return optionext.None[T]()
	}
	v := i.slice[0]
	i.slice = i.slice[1:]
	return optionext.Some(v)
}

// Iter is a convenience function that converts the slice iterator into an `*Iterate[T]`.
func (i *sliceIterator[T]) Iter() *Iterate[T] {
	return Iter[T](i)
}
