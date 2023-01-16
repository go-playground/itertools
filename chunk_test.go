package itertools

import (
	"testing"
)

func TestChunk(t *testing.T) {
	iter := WrapSlice([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).Iter().Chunk(4)
	chunk1 := iter.Next()
	assertEqual(t, chunk1.IsSome(), true)
	slice1 := chunk1.Unwrap()
	assertEqual(t, len(slice1), 4)
	assertEqual(t, slice1[0], 0)
	assertEqual(t, slice1[1], 1)
	assertEqual(t, slice1[2], 2)
	assertEqual(t, slice1[3], 3)

	chunk2 := iter.Next()
	assertEqual(t, chunk2.IsSome(), true)
	slice2 := chunk2.Unwrap()
	assertEqual(t, len(slice2), 4)
	assertEqual(t, slice2[0], 4)
	assertEqual(t, slice2[1], 5)
	assertEqual(t, slice2[2], 6)
	assertEqual(t, slice2[3], 7)

	chunk3 := iter.Next()
	assertEqual(t, chunk3.IsSome(), true)
	slice3 := chunk3.Unwrap()
	assertEqual(t, len(slice3), 2)
	assertEqual(t, slice3[0], 8)
	assertEqual(t, slice3[1], 9)

	assertEqual(t, iter.Next().IsNone(), true)
}
