package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// StepBy returns a `StepByIterator[T]` for use.
func StepBy[T any](iterator Iterator[T], step int) *StepByIterator[T] {
	return &StepByIterator[T]{
		iterator: iterator,
		step:     step,
		first:    true,
	}
}

// StepByIterator is an iterator starting at the same point, but stepping by the given amount at each iteration.
//
// The first element is always returned before the stepping begins.
type StepByIterator[T any] struct {
	iterator Iterator[T]
	step     int
	first    bool
}

// Next returns the next element advancing by the provided step or end of iterator and will ignore errors
// returned from the elements being stepped over.
func (i *StepByIterator[T]) Next() optionext.Option[T] {
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
