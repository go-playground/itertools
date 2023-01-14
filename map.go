package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"sync"
)

// Entry represents a single Map entry.
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// WrapMap creates a new iterator for transformation of types.
func WrapMap[K comparable, V any](m map[K]V) *mapWrapper[K, V] {
	return &mapWrapper[K, V]{
		mp: m,
	}
}

// mapWrapper is used to transform elements from one type to another.
type mapWrapper[K comparable, V any] struct {
	mp       map[K]V
	m        sync.Mutex
	parallel bool
}

// Next returns the next transformed element or None if at the end of the iterator.
//
// Warning: This consumes(removes) the map entries as it iterates.
func (i *mapWrapper[K, V]) Next() optionext.Option[Entry[K, V]] {
	if i.parallel {
		i.m.Lock()
		defer i.m.Unlock()
	}
	for k, v := range i.mp {
		delete(i.mp, k)
		return optionext.Some(Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return optionext.None[Entry[K, V]]()
}

// Iter is a convenience function that converts the map iterator into an `*Iterate[T]`.
func (i *mapWrapper[K, V]) Iter() *Iterate[Entry[K, V], struct{}] {
	return IterMap[Entry[K, V], struct{}](i)
}

// IterPar is a convenience function that converts the map iterator into a parallel `*Iterate[T]`.
//
// This causes the Next function to return elements protected by a Mutex.
func (i *mapWrapper[K, V]) IterPar() *Iterate[Entry[K, V], struct{}] {
	i.parallel = true
	return IterMapPar[Entry[K, V], struct{}](i)
}

// Retain retains only the elements specified by the function and removes others.
func (i *mapWrapper[K, V]) Retain(fn func(k K, v V) bool) *mapWrapper[K, V] {
	for k, v := range i.mp {
		if fn(k, v) {
			continue
		}
		delete(i.mp, k)
	}
	return i
}

// Len returns the underlying maps length.
func (i *mapWrapper[K, V]) Len() int {
	return len(i.mp)
}
