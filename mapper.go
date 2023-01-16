package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Map creates a new iterator for transformation of types.
func Map[T any, I Iterator[T], MAP any](iterator I, fn MapFn[T, MAP]) mapper[T, I, MAP] {
	return mapper[T, I, MAP]{
		iterator: iterator,
		fn:       fn,
	}
}

// MapFn represents the mapWrapper transformation function.
type MapFn[T, MAP any] func(v T) MAP

// mapWrapper is used to transform elements from one type to another.
type mapper[T any, I Iterator[T], MAP any] struct {
	iterator I
	fn       MapFn[T, MAP]
}

// Next returns the next transformed element or None if at the end of the iterator.
func (i mapper[T, I, MAP]) Next() optionext.Option[MAP] {
	v := i.iterator.Next()
	if v.IsNone() {
		return optionext.None[MAP]()
	}
	return optionext.Some(i.fn(v.Unwrap()))
}

// Iter is a convenience function that converts the map iterator into an `*Iterate[T]`.
func (i mapper[T, I, MAP]) Iter() Iterate[MAP, Iterator[MAP], struct{}] {
	return Iter[MAP, Iterator[MAP]](i)
}
