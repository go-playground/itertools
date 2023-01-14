package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"runtime"
	"sync"
	"sync/atomic"
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

// IterPar creates a new iterator with helper functions that spins up goroutines to execute the
// Collect, CollectIter, Count and ForEach functions concurrently for the number of CPUs.
//
// The results/execution order is not guaranteed when running in parallel.
//
// It defaults the Map() function to struct{}. Use IterMap() if you wish to specify a type.
func IterPar[T any](iterator Iterator[T]) *Iterate[T, struct{}] {
	return &Iterate[T, struct{}]{
		iterator: iterator,
		cpus:     runtime.NumCPU(),
	}
}

// IterMapPar creates a new iterator with helper that spins up goroutines to execute the
// Collect, CollectIter, Count and ForEach functions concurrently for the number of CPUs.
//
// The results/execution order is not guaranteed when running in parallel.
//
// It accepts a map type `V` to allow for usage of the `Map` and `CollectMap` helper function inline.
// You must use the Map() function standalone otherwise.
func IterMapPar[T, V any](iterator Iterator[T]) *Iterate[T, V] {
	return &Iterate[T, V]{
		iterator: iterator,
		cpus:     runtime.NumCPU(),
	}
}

// Iterate is an iterator with attached helper functions
type Iterate[T, V any] struct {
	iterator Iterator[T]
	cpus     int
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

// Find searches for the next element of an iterator that satisfies the function.
func (i *Iterate[T, V]) Find(fn func(T) bool) (result optionext.Option[T]) {
	for {
		result = i.iterator.Next()
		if result.IsNone() || fn(result.Unwrap()) {
			return
		}
	}
}

// All returns true if all element matches the function return, false otherwise.
//
// This will run in parallel is using a parallel iterator.
func (i *Iterate[T, V]) All(fn func(T) bool) (isAll bool) {
	if i.cpus > 0 {
		var b atomic.Bool
		b.Store(true)
		i.forEach(func(v T) (stop bool) {
			if fn(v) {
				return false
			}
			b.Store(false)
			return true
		})
		return b.Load()
	} else {
		var checked bool
		i.forEach(func(v T) (stop bool) {
			checked = fn(v)
			return !checked
		})
		return checked
	}
}

// Any returns true if any element matches the function return, false otherwise.
//
// This will run in parallel is using a parallel iterator.
func (i *Iterate[T, V]) Any(fn func(T) bool) (isAny bool) {
	if i.cpus > 0 {
		var b atomic.Bool
		i.forEach(func(v T) (stop bool) {
			match := fn(v)
			if match {
				b.Store(true)
			}
			return match
		})
		return b.Load()
	} else {
		i.forEach(func(v T) (stop bool) {
			isAny = fn(v)
			return isAny
		})
		return
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
//
// This will run in parallel is using a parallel iterator.
func (i *Iterate[T, V]) Count() int {
	if i.cpus > 0 {
		var j int64
		i.ForEach(func(_ T) {
			atomic.AddInt64(&j, 1)
		})
		return int(j)
	} else {
		var j int
		i.ForEach(func(_ T) {
			j++
		})
		return j
	}
}

// Collect transforms an iterator into a sliceWrapper.
//
// This will run in parallel is using a parallel iterator.
func (i *Iterate[T, V]) Collect() (results []T) {
	if i.cpus > 0 {
		out := make(chan T, i.cpus)
		go func() {
			i.ForEach(func(v T) {
				out <- v
			})
			close(out)
		}()
		for v := range out {
			results = append(results, v)
		}
	} else {
		i.ForEach(func(v T) {
			results = append(results, v)
		})
	}
	return
}

// ForEach runs the provided function for each element until completion.
//
// This will run in parallel is using a parallel iterator.
func (i *Iterate[T, V]) ForEach(fn func(T)) {
	i.forEach(func(t T) (stop bool) {
		fn(t)
		return false
	})
}

// forEach is an early cancellable form of `ForEach`
func (i *Iterate[T, V]) forEach(fn func(T) (stop bool)) {
	if i.cpus > 0 {
		stopEarly := make(chan struct{})
		var stopOnce sync.Once
		wg := new(sync.WaitGroup)
		for j := 0; j < i.cpus; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
			FOR:
				for {
					select {
					case <-stopEarly:
						break FOR
					default:
						v := i.iterator.Next()
						if v.IsNone() || fn(v.Unwrap()) {
							stopOnce.Do(func() {
								close(stopEarly)
							})
							break FOR
						}

					}
				}
			}()
		}
		wg.Wait()
	} else {
		for {
			v := i.iterator.Next()
			if v.IsNone() || fn(v.Unwrap()) {
				break
			}
		}
	}
}

// CollectIter transforms an iterator into a sliceWrapper and returns a *sliceWrapper in order to
// run additional functions inline such as Sort().
//
// eg. .Filter(...).CollectIter().Sort(...).WrapSlice()
//
// This will run in parallel is using a parallel iterator.
func (i *Iterate[T, V]) CollectIter() *sliceWrapper[T, struct{}] {
	return WrapSlice[T](i.Collect())
}

// Peekable returns a `PeekableIterator[T]` that wraps the current iterator.
//
// NOTE: Peekable iterators are commonly the LAST in a chain of iterators.
func (i *Iterate[T, V]) Peekable() *PeekableIterator[T] {
	return Peekable[T](i.iterator)
}
