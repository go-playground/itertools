package itertools

import optionext "github.com/go-playground/pkg/v5/values/option"

// Chunk creates a new `Chunker` for use.
//
// The default Map type is struct{}, see `ChunkWithMap` for details.
func Chunk[T any, I Iterator[T]](iterator I, size int) Chunker[T, I, struct{}] {
	return ChunkWithMap[T, I, struct{}](iterator, size)
}

// ChunkWithMap creates a new `Chunker` for use that accepts a Map type for use with `Iterate`.
func ChunkWithMap[T any, I Iterator[T], MAP any](iterator I, size int) Chunker[T, I, MAP] {
	return Chunker[T, I, MAP]{
		iterator: iterator,
		size:     size,
	}
}

// Chunker chunks the returned elements into slices of specified size.
//
// The last returned slice is NOT guaranteed to be of exact size unless the input exactly lines up.
type Chunker[T any, I Iterator[T], MAP any] struct {
	iterator I
	size     int
}

// Next yields the next set of elements from the iterator.
func (i Chunker[T, I, MAP]) Next() optionext.Option[[]T] {
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

//// Wish this was possible but the Go Compiler sees this as infinite recursion and it looks like nobody's interested in
//// fixing that :( https://github.com/golang/go/issues/50215 That's OK it works perfectly fine in Rust :P
////
//// Instead you must separate your inline statements eg. iter := WrapSlice([]int{0}); Iter(Chunk(iter).Filter().....
////
//// // Iter is a convenience function that converts the `chunker` iterator into an `*Iterate[T]`.
//func (i *Chunker[T, I, MAP]) Iter() Iterate[[]T, Iterator[[]T], MAP] {
//	return IterMap[[]T, Iterator[[]T], MAP](i)
//}
