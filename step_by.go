package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// StepBy returns a `stepByIterator[T]` for use.
func StepBy[T any, I Iterator[T]](iterator I, step int) *stepByIterator[T, I, struct{}] {
	return StepByWithMap[T, I, struct{}](iterator, step)
}

// StepByWithMap returns a `stepByIterator[T]` for use and can specify a future `Map` type conversion.
func StepByWithMap[T any, I Iterator[T], MAP any](iterator I, step int) *stepByIterator[T, I, MAP] {
	return &stepByIterator[T, I, MAP]{
		iterator: iterator,
		step:     step,
		first:    true,
	}
}

// stepByIterator is an iterator starting at the same point, but stepping by the given amount at each iteration.
//
// The first element is always returned before the stepping begins.
type stepByIterator[T any, I Iterator[T], MAP any] struct {
	iterator I
	step     int
	first    bool
}

// Next returns the next element advancing by the provided step or end of iterator and will ignore errors
// returned from the elements being stepped over.
func (i *stepByIterator[T, I, MAP]) Next() optionext.Option[T] {
	if i.first {
		i.first = false
		return i.iterator.Next()
	}
	var v optionext.Option[T]
	for j := 0; j < i.step; j++ {
		v = i.iterator.Next()
		if v.IsNone() {
			return v
		}
	}
	return v
}

// Iter is a convenience function that converts the `stepByIterator` iterator into an `*Iterate[T]`.
func (i *stepByIterator[T, I, MAP]) Iter() Iterate[T, Iterator[T], MAP] {
	return IterMap[T, Iterator[T], MAP](i)
}
