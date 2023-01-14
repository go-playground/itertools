package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Map creates a new iterator for transformation of types.
func Map[T, V any](iterator Iterator[T], fn MapFn[T, V]) *mapper[T, V] {
	return &mapper[T, V]{
		iterator: iterator,
		fn:       fn,
	}
}

// MapFn represents the mapWrapper transformation function.
type MapFn[T, V any] func(v T) V

// mapWrapper is used to transform elements from one type to another.
type mapper[T, V any] struct {
	iterator Iterator[T]
	fn       MapFn[T, V]
}

// Next returns the next transformed element or None if at the end of the iterator.
func (i *mapper[T, V]) Next() optionext.Option[V] {
	v := i.iterator.Next()
	if v.IsNone() {
		return optionext.None[V]()
	}
	return optionext.Some(i.fn(v.Unwrap()))
}

// Iter is a convenience function that converts the map iterator into an `*Iterate[T]`.
func (i *mapper[T, V]) Iter() *Iterate[V, struct{}] {
	return Iter[V](i)
}

// IterPar is a convenience function that converts the map iterator into a parallel `*Iterate[T]`.
func (i *mapper[T, V]) IterPar() *Iterate[V, struct{}] {
	return IterPar[V](i)
}
