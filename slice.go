package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"sort"
)

// SliceIter accepts and turns a slice into an iterator.
//
// This default the Map type to struct{} when none is required. See SliceIterMap if one is needed.
func SliceIter[T any](slice []T) *sliceIterator[T, struct{}] {
	return &sliceIterator[T, struct{}]{
		slice: slice,
	}
}

// SliceIterMap accepts and turns a slice into an iterator with a map type specified for Iter() to allow the Map helper
// function.
func SliceIterMap[T, V any](slice []T) *sliceIterator[T, V] {
	return &sliceIterator[T, V]{
		slice: slice,
	}
}

type sliceIterator[T, V any] struct {
	slice []T
}

func (i *sliceIterator[T, V]) Next() optionext.Option[T] {
	if len(i.slice) == 0 {
		return optionext.None[T]()
	}
	v := i.slice[0]
	i.slice = i.slice[1:]
	return optionext.Some(v)
}

// Iter is a convenience function that converts the slice iterator into an `*Iterate[T]`.
func (i *sliceIterator[T, V]) Iter() *Iterate[T, V] {
	return IterMap[T, V](i)
}

// Slice returns the underlying slice wrapped by the *sliceIterator.
func (i *sliceIterator[T, V]) Slice() []T {
	return i.slice
}

// Len returns the length of the underlying slice.
func (i *sliceIterator[T, V]) Len() int {
	return len(i.slice)
}

// Cap returns the capacity of the underlying slice.
func (i *sliceIterator[T, V]) Cap() int {
	return cap(i.slice)
}

// Sort sorts the slice x given the provided less function.
//
// The sort is not guaranteed to be stable: equal elements
// may be reversed from their original order.
// For a stable sort, use SortStable.
//
// `T` must be comparable.
func (i *sliceIterator[T, V]) Sort(less func(i T, j T) bool) *sliceIterator[T, V] {
	sort.Slice(i.slice, func(j, k int) bool {
		return less(i.slice[j], i.slice[k])
	})
	return i
}

// SortStable sorts the slice x using the provided less
// function, keeping equal elements in their original order.
func (i *sliceIterator[T, V]) SortStable(less func(i T, j T) bool) *sliceIterator[T, V] {
	sort.SliceStable(i.slice, func(j, k int) bool {
		return less(i.slice[j], i.slice[k])
	})
	return i
}

// Retain retains only the elements specified by the function.
//
// This shuffles and resizes the underlying slice.
func (i *sliceIterator[T, V]) Retain(fn func(v T) bool) *sliceIterator[T, V] {
	var j int
	for _, v := range i.slice {
		if fn(v) {
			i.slice[j] = v
			j++
		}
	}
	i.slice = i.slice[:j]
	return i
}
