package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"sort"
)

// WrapSlice accepts and turns a sliceWrapper into an iterator.
//
// The default the Map type to struct{} when none is required. See WrapSliceMap if one is needed.
func WrapSlice[T any](slice []T) *sliceWrapper[T, struct{}] {
	return WrapSliceMap[T, struct{}](slice)
}

// WrapSliceMap accepts and turns a sliceWrapper into an iterator with a map type specified for IterPar() to allow the Map helper
// function.
func WrapSliceMap[T, V any](slice []T) *sliceWrapper[T, V] {
	return &sliceWrapper[T, V]{
		slice: slice,
	}
}

type sliceWrapper[T, V any] struct {
	slice []T
}

func (i *sliceWrapper[T, V]) Next() optionext.Option[T] {
	if len(i.slice) == 0 {
		return optionext.None[T]()
	}
	v := i.slice[0]
	i.slice = i.slice[1:]
	return optionext.Some(v)
}

// Iter is a convenience function that converts the sliceWrapper iterator into an `*Iterate[T]`.
func (i *sliceWrapper[T, V]) Iter() *Iterate[T, V] {
	return IterMap[T, V](i)
}

// Slice returns the underlying sliceWrapper wrapped by the *sliceWrapper.
func (i *sliceWrapper[T, V]) Slice() []T {
	return i.slice
}

// Len returns the length of the underlying sliceWrapper.
func (i *sliceWrapper[T, V]) Len() int {
	return len(i.slice)
}

// Cap returns the capacity of the underlying sliceWrapper.
func (i *sliceWrapper[T, V]) Cap() int {
	return cap(i.slice)
}

// Sort sorts the sliceWrapper x given the provided less function.
//
// The sort is not guaranteed to be stable: equal elements
// may be reversed from their original order.
// For a stable sort, use SortStable.
//
// `T` must be comparable.
func (i *sliceWrapper[T, V]) Sort(less func(i T, j T) bool) *sliceWrapper[T, V] {
	sort.Slice(i.slice, func(j, k int) bool {
		return less(i.slice[j], i.slice[k])
	})
	return i
}

// SortStable sorts the sliceWrapper x using the provided less
// function, keeping equal elements in their original order.
func (i *sliceWrapper[T, V]) SortStable(less func(i T, j T) bool) *sliceWrapper[T, V] {
	sort.SliceStable(i.slice, func(j, k int) bool {
		return less(i.slice[j], i.slice[k])
	})
	return i
}

// Retain retains only the elements specified by the function.
//
// This shuffles and resizes the underlying sliceWrapper.
func (i *sliceWrapper[T, V]) Retain(fn func(v T) bool) *sliceWrapper[T, V] {
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
