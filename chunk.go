package itertools

import optionext "github.com/go-playground/pkg/v5/values/option"

// Chunk creates a new `Chunker` for use.
//
// The default Map type is struct{}, see `ChunkMap` for details.
func Chunk[T any](iterator Iterator[T], size int) *Chunker[T, struct{}] {
	return ChunkMap[T, struct{}](iterator, size)
}

// ChunkMap creates a new `Chunker` for use that accepts a Map type for use with `Iterate`.
func ChunkMap[T, MAP any](iterator Iterator[T], size int) *Chunker[T, MAP] {
	return &Chunker[T, MAP]{
		iterator: iterator,
		size:     size,
	}
}

// Chunker chunks the returned elements into slices of specified size.
//
// The last returned slice is NOT guaranteed to be of exact size unless the input exactly lines up.
type Chunker[T, MAP any] struct {
	iterator Iterator[T]
	size     int
}

// Next yields the next set of elements from the iterator.
func (i *Chunker[T, MAP]) Next() optionext.Option[[]T] {
	chunk := make([]T, 0, i.size)
	for {
		v := i.iterator.Next()
		if v.IsNone() {
			break
		}
		chunk = append(chunk, v.Unwrap())
		if len(chunk) == cap(chunk) {
			break
		}
	}
	if len(chunk) == 0 {
		return optionext.None[[]T]()
	}
	return optionext.Some(chunk)
}

// Wish this was possible but the Go Compiler sees this as infinite recursion and it looks like nobody's interested in
// fixing that :( https://github.com/golang/go/issues/50215 That's OK it works perfectly fine in Rust :P
//
// Instead you must separate your inline statements eg. iter := WrapSlice([]int{0}); Iter(Chunk(iter).Filter().....
//
//// Iter is a convenience function that converts the `chunker` iterator into an `*Iterate[T]`.
//func (i *Chunker[T, MAP]) Iter() *Iterate[[]T, MAP] {
//	return IterMap[[]T, MAP](i)
//}
