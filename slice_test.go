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
	iter = Filter(WrapSlice(slice).IntoIter(), func(v int) bool {
		return v < 9
	}).Iter().CollectIter()
	assertEqual(t, iter.Next(), optionext.Some(9))
	assertEqual(t, iter.Next(), optionext.None[int]())

	// Test Iter Filter
	iter = WrapSlice(slice).Iter().Filter(func(v int) bool {
		return v < 9
	}).CollectIter()
	assertEqual(t, iter.Next(), optionext.Some(9))
	assertEqual(t, iter.Next(), optionext.None[int]())

	// Test Retain
	iter = WrapSlice(slice).Retain(func(v int) bool {
		return v == 3
	})
	assertEqual(t, iter.Len(), 1)
	assertEqual(t, iter.Next(), optionext.Some(3))
	assertEqual(t, iter.Next(), optionext.None[int]())

	// Test Retain
	slice = RetainSlice([]int{0, 1, 2, 3}, func(v int) bool {
		return v == 2
	})
	assertEqual(t, len(slice), 1)
	assertEqual(t, slice[0], 2)

	// Test sort
	slice = []int{0, 1, 2, 3}
	iterChain := WrapSlice(slice).Iter().Chain(WrapSlice(slice).IntoIter())
	assertEqual(t, iterChain.Next(), optionext.Some(0))
	assertEqual(t, iterChain.Next(), optionext.Some(1))
	assertEqual(t, iterChain.Next(), optionext.Some(2))
	assertEqual(t, iterChain.Next(), optionext.Some(3))
	assertEqual(t, iterChain.Next(), optionext.Some(0))
	assertEqual(t, iterChain.Next(), optionext.Some(1))
	assertEqual(t, iterChain.Next(), optionext.Some(2))
	assertEqual(t, iterChain.Next(), optionext.Some(3))
	assertEqual(t, iterChain.Next(), optionext.None[int]())

	iterChain = Chain[int](WrapSlice(slice).IntoIter(), WrapSlice(slice).IntoIter()).Iter()
	assertEqual(t, iterChain.Next(), optionext.Some(0))
	assertEqual(t, iterChain.Next(), optionext.Some(1))
	assertEqual(t, iterChain.Next(), optionext.Some(2))
	assertEqual(t, iterChain.Next(), optionext.Some(3))
	assertEqual(t, iterChain.Next(), optionext.Some(0))
	assertEqual(t, iterChain.Next(), optionext.Some(1))
	assertEqual(t, iterChain.Next(), optionext.Some(2))
	assertEqual(t, iterChain.Next(), optionext.Some(3))
	assertEqual(t, iterChain.Next(), optionext.None[int]())
}

func assertEqual[T comparable](t *testing.T, l, r T) {
	if l != r {
		t.Fatalf("expected `%#v` to equal `%#v`", l, r)
	}
}

func stdRetain(s []int) []int {
	var j int
	for _, v := range s {
		if v == 1 {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

func stdRetainFn(s []int, fn func(v int) bool) []int {
	var j int
	for _, v := range s {
		if fn(v) {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

func BenchmarkSTDRetain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stdRetain(makeSlice())
	}
}

func BenchmarkSTDFnRetain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stdRetainFn(makeSlice(), func(v int) bool {
			return v == 1
		})
	}
}

func BenchmarkSliceWrapper_Retain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WrapSlice(makeSlice()).Retain(func(v int) bool {
			return v == 1
		})
	}
}

func BenchmarkRetainSlice_Retain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RetainSlice(makeSlice(), func(v int) bool {
			return v == 1
		})
	}
}
