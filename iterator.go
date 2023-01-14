package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Iterator is an interface representing something that iterates using the Next method.
type Iterator[T any] interface {
	// Next advances the iterator and returns the next value.
	//
	// Returns an Option with value of None when iteration has finished.
	Next() optionext.Option[T]
}

// Iter creates a new iterator with helper functions.
//
// It defaults the Map() function to struct{}. Use IterMap() if you wish to specify a type.
func Iter[T any](iterator Iterator[T]) *Iterate[T, struct{}] {
	return &Iterate[T, struct{}]{
		iterator: iterator,
	}
}

// IterMap creates a new iterator with helper functions.
//
// It accepts a map type `V` to allow for usage of the `Map` and `CollectMap` helper function inline.
// You must use the Map() function standalone otherwise.
func IterMap[T, V any](iterator Iterator[T]) *Iterate[T, V] {
	return &Iterate[T, V]{
		iterator: iterator,
	}
}

// Iterate is an iterator with attached helper functions
type Iterate[T, V any] struct {
	iterator Iterator[T]
}

// Next returns the new iterator value
func (i *Iterate[T, V]) Next() optionext.Option[T] {
	return i.iterator.Next()
}

// Map accepts a `FilterFn[T]` to filter items.
func (i *Iterate[T, V]) Map(fn MapFn[T, V]) *mapper[T, V] {
	return Map[T, V](i.iterator, fn)
}

// Filter accepts a `FilterFn[T]` to filter items.
func (i *Iterate[T, V]) Filter(fn FilterFn[T]) *Iterate[T, V] {
	i.iterator = Filter[T](i.iterator, fn)
	return i
}

// Chain creates a new ChainIterator for use.
func (i *Iterate[T, V]) Chain(iterator Iterator[T]) *Iterate[T, V] {
	i.iterator = Chain[T](i.iterator, iterator)
	return i
}

// Take yields elements until n elements are yielded or the end of the iterator is reached (whichever happens first)
func (i *Iterate[T, V]) Take(n int) *Iterate[T, V] {
	i.iterator = Take[T](i.iterator, n)
	return i
}

// TakeWhile yields elements while the function return true or the end of the iterator is reached (whichever happens first)
func (i *Iterate[T, V]) TakeWhile(fn TakeWhileFn[T]) *Iterate[T, V] {
	i.iterator = TakeWhile[T](i.iterator, fn)
	return i
}

// StepBy returns a `StepByIterator[T, I]` starting at the same point, but stepping by the given amount at each iteration.
//
// The first element is always returned before the stepping begins.
func (i *Iterate[T, V]) StepBy(step int) *Iterate[T, V] {
	i.iterator = StepBy[T](i.iterator, step)
	return i
}

// Find searches for an element of an iterator that satisfies the function.
func (i *Iterate[T, V]) Find(fn func(T) bool) (result optionext.Option[T]) {
	for {
		result = i.iterator.Next()
		if result.IsNone() || fn(result.Unwrap()) {
			return
		}
	}
}

// All returns true if all element matches the function return, false otherwise.
func (i *Iterate[T, V]) All(fn func(T) bool) bool {
	var checked bool
	for {
		result := i.iterator.Next()
		if result.IsNone() {
			return checked
		} else if !fn(result.Unwrap()) {
			return false
		}
		checked = true
	}
}

// Any returns true if any element matches the function return, false otherwise.
func (i *Iterate[T, V]) Any(fn func(T) bool) bool {
	for {
		result := i.iterator.Next()
		if result.IsNone() {
			return false
		} else if fn(result.Unwrap()) {
			return true
		}
	}
}

// Position searches for an element in an iterator, returning its index.
func (i *Iterate[T, V]) Position(fn func(T) bool) optionext.Option[int] {
	var j int
	for {
		result := i.iterator.Next()
		if result.IsNone() {
			return optionext.None[int]()
		} else if fn(result.Unwrap()) {
			return optionext.Some(j)
		}
		j++
	}
}

// Count consumes the iterator and returns count if iterations.
func (i *Iterate[T, V]) Count() int {
	var j int
	for {
		result := i.iterator.Next()
		if result.IsNone() {
			return j
		}
		j++
	}
}

// Collect transforms an iterator into a slice.
func (i *Iterate[T, V]) Collect() (results []T) {
	for {
		v := i.iterator.Next()
		if v.IsNone() {
			return results
		}
		results = append(results, v.Unwrap())
	}
}

// CollectIter transforms an iterator into a slice and returns a *sliceIterator in order to
// run additional functions inline such as Sort().
//
// eg. .Filter(...).CollectIter().Sort(...).Slice()
func (i *Iterate[T, V]) CollectIter() *sliceIterator[T, struct{}] {
	return SliceIter[T](i.Collect())
}

// Peekable returns a `PeekableIterator[T]` that wraps the current iterator.
//
// NOTE: Peekable iterators are commonly the LAST in a chain of iterators.
func (i *Iterate[T, V]) Peekable() *PeekableIterator[T] {
	return Peekable[T](i.iterator)
}
