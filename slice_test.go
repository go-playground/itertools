package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"testing"
)

func TestSlice(t *testing.T) {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Test WrapSlice, Len, Cap
	iter := WrapSlice(slice)
	assertEqual(t, iter.Len(), 10)
	assertEqual(t, iter.Cap(), 10)
	assertEqual(t, len(iter.Slice()), 10)

	// Test Next
	iter = WrapSlice(slice)
	assertEqual(t, iter.Next(), optionext.Some(0))
	assertEqual(t, iter.Next(), optionext.Some(1))
	assertEqual(t, iter.Next(), optionext.Some(2))
	assertEqual(t, iter.Next(), optionext.Some(3))
	assertEqual(t, iter.Next(), optionext.Some(4))
	assertEqual(t, iter.Next(), optionext.Some(5))
	assertEqual(t, iter.Next(), optionext.Some(6))
	assertEqual(t, iter.Next(), optionext.Some(7))
	assertEqual(t, iter.Next(), optionext.Some(8))
	assertEqual(t, iter.Next(), optionext.Some(9))
	assertEqual(t, iter.Next(), optionext.None[int]())

	// Test sort
	iter = WrapSlice(slice).Sort(func(i int, j int) bool {
		return i > j
	})
	assertEqual(t, iter.Next(), optionext.Some(9))
	assertEqual(t, iter.Next(), optionext.Some(8))
	assertEqual(t, iter.Next(), optionext.Some(7))
	assertEqual(t, iter.Next(), optionext.Some(6))
	assertEqual(t, iter.Next(), optionext.Some(5))
	assertEqual(t, iter.Next(), optionext.Some(4))
	assertEqual(t, iter.Next(), optionext.Some(3))
	assertEqual(t, iter.Next(), optionext.Some(2))
	assertEqual(t, iter.Next(), optionext.Some(1))
	assertEqual(t, iter.Next(), optionext.Some(0))
	assertEqual(t, iter.Next(), optionext.None[int]())

	// Test sort stable
	iter = WrapSlice(slice).SortStable(func(i int, j int) bool {
		return i > j
	})
	assertEqual(t, iter.Next(), optionext.Some(9))
	assertEqual(t, iter.Next(), optionext.Some(8))
	assertEqual(t, iter.Next(), optionext.Some(7))
	assertEqual(t, iter.Next(), optionext.Some(6))
	assertEqual(t, iter.Next(), optionext.Some(5))
	assertEqual(t, iter.Next(), optionext.Some(4))
	assertEqual(t, iter.Next(), optionext.Some(3))
	assertEqual(t, iter.Next(), optionext.Some(2))
	assertEqual(t, iter.Next(), optionext.Some(1))
	assertEqual(t, iter.Next(), optionext.Some(0))
	assertEqual(t, iter.Next(), optionext.None[int]())

	// Test Iter Filter
	iter = WrapSlice(slice).Iter().Filter(func(v int) bool {
		return v < 9
	}).CollectIter()
	assertEqual(t, iter.Next(), optionext.Some(9))
	assertEqual(t, iter.Next(), optionext.None[int]())

	// Test IterPar Filter
	iter = WrapSlice(slice).IterPar().Filter(func(v int) bool {
		return v < 9
	}).CollectIter()
	assertEqual(t, iter.Next(), optionext.Some(9))
	assertEqual(t, iter.Next(), optionext.None[int]())

	// Test Retain
	iter = WrapSlice(slice).Retain(func(v int) bool {
		return v == 3
	})
	assertEqual(t, iter.Next(), optionext.Some(3))
	assertEqual(t, iter.Next(), optionext.None[int]())
}

func assertEqual[T comparable](t *testing.T, l, r T) {
	if l != r {
		t.Fatalf("expected `%#v` to equal `%#v`", l, r)
	}
}
