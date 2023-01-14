package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
)

// Entry represents a single Map entry.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// MapIter creates a new iterator for transformation of types.
func MapIter[K comparable, V any](m map[K]V) *mapIterator[K, V] {
	return &mapIterator[K, V]{
		m: m,
	}
}

// mapIterator is used to transform elements from one type to another.
type mapIterator[K comparable, V any] struct {
	m map[K]V
}

// Next returns the next transformed element or None if at the end of the iterator.
//
// Warning: This consumes(removes) the map entries as it iterates.
func (i *mapIterator[K, V]) Next() optionext.Option[Entry[K, V]] {
	for k, v := range i.m {
		delete(i.m, k)
		return optionext.Some(Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return optionext.None[Entry[K, V]]()
}

// Iter is a convenience function that converts the map iterator into an `*Iterate[T]`.
func (i *mapIterator[K, V]) Iter() *Iterate[Entry[K, V], struct{}] {
	return IterMap[Entry[K, V], struct{}](i)
}

// Retain retains only the elements specified by the function and removes others.
func (i *mapIterator[K, V]) Retain(fn func(k K, v V) bool) *mapIterator[K, V] {
	for k, v := range i.m {
		if fn(k, v) {
			continue
		}
		delete(i.m, k)
	}
	return i
}

// TODO: Add Len function
