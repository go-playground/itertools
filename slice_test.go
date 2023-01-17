package itertools

import (
	. "github.com/go-playground/assert/v2"
	sliceext "github.com/go-playground/pkg/v5/slice"
	optionext "github.com/go-playground/pkg/v5/values/option"
	"strconv"
	"testing"
)

func TestSlice(t *testing.T) {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	// Test WrapSlice, Len, Cap
	iter := WrapSlice(slice)
	Equal(t, iter.Len(), 10)
	Equal(t, iter.Cap(), 10)
	Equal(t, len(iter.Slice()), 10)

	// Test Next
	iter = WrapSlice(slice)
	Equal(t, iter.Next(), optionext.Some(0))
	Equal(t, iter.Next(), optionext.Some(1))
	Equal(t, iter.Next(), optionext.Some(2))
	Equal(t, iter.Next(), optionext.Some(3))
	Equal(t, iter.Next(), optionext.Some(4))
	Equal(t, iter.Next(), optionext.Some(5))
	Equal(t, iter.Next(), optionext.Some(6))
	Equal(t, iter.Next(), optionext.Some(7))
	Equal(t, iter.Next(), optionext.Some(8))
	Equal(t, iter.Next(), optionext.Some(9))
	Equal(t, iter.Next(), optionext.None[int]())

	// Test sort
	iter = WrapSlice(slice).Sort(func(i int, j int) bool {
		return i > j
	})
	Equal(t, iter.Next(), optionext.Some(9))
	Equal(t, iter.Next(), optionext.Some(8))
	Equal(t, iter.Next(), optionext.Some(7))
	Equal(t, iter.Next(), optionext.Some(6))
	Equal(t, iter.Next(), optionext.Some(5))
	Equal(t, iter.Next(), optionext.Some(4))
	Equal(t, iter.Next(), optionext.Some(3))
	Equal(t, iter.Next(), optionext.Some(2))
	Equal(t, iter.Next(), optionext.Some(1))
	Equal(t, iter.Next(), optionext.Some(0))
	Equal(t, iter.Next(), optionext.None[int]())

	// Test sort stable
	iter = WrapSlice(slice).SortStable(func(i int, j int) bool {
		return i > j
	})
	Equal(t, iter.Next(), optionext.Some(9))
	Equal(t, iter.Next(), optionext.Some(8))
	Equal(t, iter.Next(), optionext.Some(7))
	Equal(t, iter.Next(), optionext.Some(6))
	Equal(t, iter.Next(), optionext.Some(5))
	Equal(t, iter.Next(), optionext.Some(4))
	Equal(t, iter.Next(), optionext.Some(3))
	Equal(t, iter.Next(), optionext.Some(2))
	Equal(t, iter.Next(), optionext.Some(1))
	Equal(t, iter.Next(), optionext.Some(0))
	Equal(t, iter.Next(), optionext.None[int]())

	// Test Iter Filter
	iter = Filter(WrapSlice(slice).IntoIter(), func(v int) bool {
		return v < 9
	}).Iter().CollectIter()
	Equal(t, iter.Next(), optionext.Some(9))
	Equal(t, iter.Next(), optionext.None[int]())

	// Test Iter Filter
	iter = WrapSlice(slice).Iter().Filter(func(v int) bool {
		return v < 9
	}).CollectIter()
	Equal(t, iter.Next(), optionext.Some(9))
	Equal(t, iter.Next(), optionext.None[int]())

	// Test Retain
	iter = WrapSlice(slice).Retain(func(v int) bool {
		return v == 3
	})
	Equal(t, iter.Len(), 1)
	Equal(t, iter.Next(), optionext.Some(3))
	Equal(t, iter.Next(), optionext.None[int]())

	// Test Filter
	iter = WrapSlice([]int{0, 1, 2, 3, 4, 5, 6, 7}).Filter(func(v int) bool {
		return v != 3
	})
	Equal(t, iter.Len(), 1)
	Equal(t, iter.Next(), optionext.Some(3))
	Equal(t, iter.Next(), optionext.None[int]())

	// Test Iter sort
	slice = []int{0, 1, 2, 3}
	iterChain := WrapSlice(slice).Iter().Chain(WrapSlice(slice).IntoIter())
	Equal(t, iterChain.Next(), optionext.Some(0))
	Equal(t, iterChain.Next(), optionext.Some(1))
	Equal(t, iterChain.Next(), optionext.Some(2))
	Equal(t, iterChain.Next(), optionext.Some(3))
	Equal(t, iterChain.Next(), optionext.Some(0))
	Equal(t, iterChain.Next(), optionext.Some(1))
	Equal(t, iterChain.Next(), optionext.Some(2))
	Equal(t, iterChain.Next(), optionext.Some(3))
	Equal(t, iterChain.Next(), optionext.None[int]())

	// Test Native sort
	slice = []int{0, 1, 2, 3}
	sorted := WrapSlice(slice).Sort(func(i int, j int) bool {
		return i > j
	})
	Equal(t, sorted.Next(), optionext.Some(3))
	Equal(t, sorted.Next(), optionext.Some(2))
	Equal(t, sorted.Next(), optionext.Some(1))
	Equal(t, sorted.Next(), optionext.Some(0))
	Equal(t, sorted.Next(), optionext.None[int]())

	slice = []int{0, 1, 2, 3}
	iterChainWrap := Chain[int](WrapSlice(slice).IntoIter(), WrapSlice(slice).IntoIter()).Iter()
	Equal(t, iterChainWrap.Next(), optionext.Some(0))
	Equal(t, iterChainWrap.Next(), optionext.Some(1))
	Equal(t, iterChainWrap.Next(), optionext.Some(2))
	Equal(t, iterChainWrap.Next(), optionext.Some(3))
	Equal(t, iterChainWrap.Next(), optionext.Some(0))
	Equal(t, iterChainWrap.Next(), optionext.Some(1))
	Equal(t, iterChainWrap.Next(), optionext.Some(2))
	Equal(t, iterChainWrap.Next(), optionext.Some(3))
	Equal(t, iterChainWrap.Next(), optionext.None[int]())

	iterMap := WrapSliceMap[int, []string]([]int{0, 1, 2, 3}).Map(make([]string, 0, 4), func(accum []string, v int) []string {
		return append(accum, strconv.Itoa(v))
	})
	Equal(t, len(iterMap), 4)
	Equal(t, iterMap[0], "0")
	Equal(t, iterMap[1], "1")
	Equal(t, iterMap[2], "2")
	Equal(t, iterMap[3], "3")
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
		sliceext.Retain(makeSlice(), func(v int) bool {
			return v == 1
		})
	}
}
