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

// Iter creates a new iterator with helper functions
func Iter[T any](iterator Iterator[T]) *Iterate[T] {
	return &Iterate[T]{
		iterator: iterator,
	}
}

// Iterate is an iterator with attached helper functions
type Iterate[T any] struct {
	iterator Iterator[T]
}

// Next returns the new iterator value
func (i *Iterate[T]) Next() optionext.Option[T] {
	return i.iterator.Next()
}

// Filter accepts a `FilterFn[T]` to filter items.
func (i *Iterate[T]) Filter(fn FilterFn[T]) *Iterate[T] {
	i.iterator = Filter[T](i.iterator, fn)
	return i
}

// Chain creates a new ChainIterator for use.
func (i *Iterate[T]) Chain(iterator Iterator[T]) *Iterate[T] {
	i.iterator = Chain[T](i.iterator, iterator)
	return i
}

// Take yields elements until n elements are yielded or the end of the iterator is reached (whichever happens first)
func (i *Iterate[T]) Take(n int) *Iterate[T] {
	i.iterator = Take[T](i.iterator, n)
	return i
}

// TakeWhile yields elements while the function return true or the end of the iterator is reached (whichever happens first)
func (i *Iterate[T]) TakeWhile(fn TakeWhileFn[T]) *Iterate[T] {
	i.iterator = TakeWhile[T](i.iterator, fn)
	return i
}

// StepBy returns a `StepByIterator[T, I]` starting at the same point, but stepping by the given amount at each iteration.
//
// The first element is always returned before the stepping begins.
func (i *Iterate[T]) StepBy(step int) *Iterate[T] {
	i.iterator = StepBy[T](i.iterator, step)
	return i
}

// Find searches for an element of an iterator that satisfies the function.
func (i *Iterate[T]) Find(fn func(T) bool) (result optionext.Option[T]) {
	for {
		result = i.iterator.Next()
		if result.IsNone() || fn(result.Unwrap()) {
			return
		}
	}
}

// All returns true if all element matches the function return, false otherwise.
func (i *Iterate[T]) All(fn func(T) bool) bool {
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
func (i *Iterate[T]) Any(fn func(T) bool) bool {
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
func (i *Iterate[T]) Position(fn func(T) bool) optionext.Option[int] {
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
func (i *Iterate[T]) Count() int {
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
func (i *Iterate[T]) Collect() (results []T) {
	for {
		v := i.iterator.Next()
		if v.IsNone() {
			return results
		}
		results = append(results, v.Unwrap())
	}
}

// Peekable returns a `PeekableIterator[T]` that wraps the current iterator.
//
// NOTE: Peekable iterators are commonly the LAST in a chain of iterators.
func (i *Iterate[T]) Peekable() *PeekableIterator[T] {
	return Peekable[T](i.iterator)
}
