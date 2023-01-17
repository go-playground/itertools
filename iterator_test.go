package itertools

import (
	. "github.com/go-playground/assert/v2"
	optionext "github.com/go-playground/pkg/v5/values/option"
	"strconv"
	"sync/atomic"
	"testing"
)

func TestIterate(t *testing.T) {

	// Test Map
	iter := WrapSliceMap[int, string](makeSlice()).Iter().Map(func(v int) string {
		return strconv.Itoa(v)
	}).Iter().CollectIter()
	Equal(t, iter.Len(), 3)
	Equal(t, iter.Next(), optionext.Some("0"))
	Equal(t, iter.Next(), optionext.Some("1"))
	Equal(t, iter.Next(), optionext.Some("2"))
	Equal(t, iter.Next(), optionext.None[string]())

	// Test Filter
	iter2 := WrapSlice(makeSlice()).Iter().Filter(func(v int) bool {
		return v != 2
	}).CollectIter()
	Equal(t, iter2.Len(), 1)
	Equal(t, iter2.Next(), optionext.Some(2))
	Equal(t, iter2.Next(), optionext.None[int]())

	// Test TakeWhile
	iter3 := WrapSlice(makeSlice()).Iter().TakeWhile(func(v int) bool {
		return v < 2
	})
	Equal(t, iter3.Next(), optionext.Some(0))
	Equal(t, iter3.Next(), optionext.Some(1))
	Equal(t, iter3.Next(), optionext.None[int]())

	// Test Take
	iter3 = WrapSlice(makeSlice()).Iter().Take(2)
	Equal(t, iter3.Next(), optionext.Some(0))
	Equal(t, iter3.Next(), optionext.Some(1))
	Equal(t, iter3.Next(), optionext.None[int]())

	// Test StepBy
	iter3 = WrapSlice(makeSlice()).Iter().StepBy(2)
	Equal(t, iter3.Next(), optionext.Some(0))
	Equal(t, iter3.Next(), optionext.Some(2))
	Equal(t, iter3.Next(), optionext.None[int]())

	// Test Find
	iter4 := WrapSlice(makeSlice()).Iter()
	Equal(t, iter4.Find(func(i int) bool {
		return i == 1
	}), optionext.Some(1))

	// Test All
	iter5 := WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.All(func(i int) bool {
		return i < 10
	}), true)
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.All(func(i int) bool {
		return i < 1
	}), false)
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.AllParallel(func(i int) bool {
		return i < 10
	}), true)
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.AllParallel(func(i int) bool {
		return i < 1
	}), false)

	// Test Any
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.Any(func(i int) bool {
		return i == 1
	}), true)
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.Any(func(i int) bool {
		return i == 10
	}), false)
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.AnyParallel(func(i int) bool {
		return i == 1
	}), true)
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.AnyParallel(func(i int) bool {
		return i == 10
	}), false)

	// Test Position
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.Position(func(i int) bool {
		return i == 1
	}), optionext.Some(1))

	// Test Count
	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.Count(), 3)

	iter5 = WrapSlice(makeSlice()).Iter()
	Equal(t, iter5.CountParallel(), 3)

	// Test ForEach
	var j int
	WrapSlice(makeSlice()).Iter().ForEach(func(_ int) {
		j++
	})
	Equal(t, j, 3)

	var k int64
	WrapSlice(makeSlice()).Iter().ForEachParallel(func(_ int) {
		atomic.AddInt64(&k, 1)
	})
	Equal(t, k, int64(3))

	// Test Chain
	iter6 := WrapSlice(makeSlice()).Iter().Chain(WrapSlice(makeSlice()).IntoIter())
	Equal(t, iter6.Next(), optionext.Some(0))
	Equal(t, iter6.Next(), optionext.Some(1))
	Equal(t, iter6.Next(), optionext.Some(2))
	Equal(t, iter6.Next(), optionext.Some(0))
	Equal(t, iter6.Next(), optionext.Some(1))
	Equal(t, iter6.Next(), optionext.Some(2))
	Equal(t, iter6.Next(), optionext.None[int]())

	// Test Peekable
	iter7 := WrapSlice(makeSlice()).Iter().Peekable()
	Equal(t, iter7.Peek(), optionext.Some(0))
	Equal(t, iter7.Next(), optionext.Some(0))
	Equal(t, iter7.Peek(), optionext.Some(1))
	Equal(t, iter7.Next(), optionext.Some(1))
	Equal(t, iter7.Peek(), optionext.Some(2))
	Equal(t, iter7.Next(), optionext.Some(2))
	Equal(t, iter7.Peek(), optionext.None[int]())
	Equal(t, iter7.Next(), optionext.None[int]())

	//	// Test Reduce
	num := WrapSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}).Iter().Reduce(func(accum int, current int) int {
		return accum + current
	})
	Equal(t, num, optionext.Some(45))

	// Test Partition
	left, right := WrapSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}).Iter().PartitionIter(func(v int) bool {
		return v%2 == 0
	})
	Equal(t, left.Next(), optionext.Some(2))
	Equal(t, left.Next(), optionext.Some(4))
	Equal(t, left.Next(), optionext.Some(6))
	Equal(t, left.Next(), optionext.Some(8))
	Equal(t, left.Next(), optionext.None[int]())
	Equal(t, right.Next(), optionext.Some(1))
	Equal(t, right.Next(), optionext.Some(3))
	Equal(t, right.Next(), optionext.Some(5))
	Equal(t, right.Next(), optionext.Some(7))
	Equal(t, right.Next(), optionext.Some(9))
	Equal(t, right.Next(), optionext.None[int]())
}

func makeSlice() []int {
	return []int{0, 1, 2}
}

func BenchmarkIterate_Complex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WrapSliceMap[int, string]([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).Iter().StepBy(2).Filter(func(v int) bool {
			return v < 6
		}).Map(func(v int) string {
			return strconv.Itoa(v)
		}).Iter().CollectIter().Sort(func(i string, j string) bool {
			return i < j
		})
	}
}
