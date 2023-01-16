package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// StepBy returns a `stepByIterator[T]` for use.
func StepBy[T any](iterator Iterator[T], step int) *stepByIterator[T, struct{}] {
	return StepByMap[T, struct{}](iterator, step)
}

// StepByMap returns a `stepByIterator[T]` for use and can specify a future `Map` type conversion.
func StepByMap[T, MAP any](iterator Iterator[T], step int) *stepByIterator[T, MAP] {
	return &stepByIterator[T, MAP]{
		iterator: iterator,
		step:     step,
		first:    true,
	}
}

// stepByIterator is an iterator starting at the same point, but stepping by the given amount at each iteration.
//
// The first element is always returned before the stepping begins.
type stepByIterator[T, MAP any] struct {
	iterator Iterator[T]
	step     int
	first    bool
}

// Next returns the next element advancing by the provided step or end of iterator and will ignore errors
// returned from the elements being stepped over.
func (i *stepByIterator[T, MAP]) Next() optionext.Option[T] {
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
func (i *stepByIterator[T, MAP]) Iter() Iterate[T, MAP] {
	return IterMap[T, MAP](i)
}
