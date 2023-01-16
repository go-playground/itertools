package itertools

import (
	. "github.com/go-playground/assert/v2"
	"testing"
)

func TestChunk(t *testing.T) {
	iter := WrapSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).Iter().Chunk(4)
	chunk1 := iter.Next()
	Equal(t, chunk1.IsSome(), true)
	slice1 := chunk1.Unwrap()
	Equal(t, len(slice1), 4)
	Equal(t, slice1[0], 0)
	Equal(t, slice1[1], 1)
	Equal(t, slice1[2], 2)
	Equal(t, slice1[3], 3)

	chunk2 := iter.Next()
	Equal(t, chunk2.IsSome(), true)
	slice2 := chunk2.Unwrap()
	Equal(t, len(slice2), 4)
	Equal(t, slice2[0], 4)
	Equal(t, slice2[1], 5)
	Equal(t, slice2[2], 6)
	Equal(t, slice2[3], 7)

	chunk3 := iter.Next()
	Equal(t, chunk3.IsSome(), true)
	slice3 := chunk3.Unwrap()
	Equal(t, len(slice3), 2)
	Equal(t, slice3[0], 8)
	Equal(t, slice3[1], 9)

	Equal(t, iter.Next().IsNone(), true)

	chunker := Chunk[int, Iterator[int]](WrapSlice([]int{1, 2, 3}).IntoIter(), 2)
	Equal(t, len(chunker.Next().Unwrap()), 2)
	Equal(t, len(chunker.Next().Unwrap()), 1)
	Equal(t, chunker.Next().IsNone(), true)
}
