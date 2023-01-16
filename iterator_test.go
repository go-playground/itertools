package itertools

import (
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
	assertEqual(t, iter.Len(), 3)
	assertEqual(t, iter.Next(), optionext.Some("0"))
	assertEqual(t, iter.Next(), optionext.Some("1"))
	assertEqual(t, iter.Next(), optionext.Some("2"))
	assertEqual(t, iter.Next(), optionext.None[string]())

	// Test Filter
	iter2 := WrapSlice(makeSlice()).Iter().Filter(func(v int) bool {
		return v != 2
	}).CollectIter()
	assertEqual(t, iter2.Len(), 1)
	assertEqual(t, iter2.Next(), optionext.Some(2))
	assertEqual(t, iter2.Next(), optionext.None[int]())

	// Test TakeWhile
	iter3 := WrapSlice(makeSlice()).Iter().TakeWhile(func(v int) bool {
		return v < 2
	})
	assertEqual(t, iter3.Next(), optionext.Some(0))
	assertEqual(t, iter3.Next(), optionext.Some(1))
	assertEqual(t, iter3.Next(), optionext.None[int]())

	// Test Take
	iter3 = WrapSlice(makeSlice()).Iter().Take(2)
	assertEqual(t, iter3.Next(), optionext.Some(0))
	assertEqual(t, iter3.Next(), optionext.Some(1))
	assertEqual(t, iter3.Next(), optionext.None[int]())

	// Test StepBy
	iter3 = WrapSlice(makeSlice()).Iter().StepBy(2)
	assertEqual(t, iter3.Next(), optionext.Some(0))
	assertEqual(t, iter3.Next(), optionext.Some(2))
	assertEqual(t, iter3.Next(), optionext.None[int]())

	// Test Find
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.Find(func(i int) bool {
		return i == 1
	}), optionext.Some(1))

	// Test All
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.All(func(i int) bool {
		return i < 10
	}), true)
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.All(func(i int) bool {
		return i < 1
	}), false)
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.AllParallel(func(i int) bool {
		return i < 10
	}), true)
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.AllParallel(func(i int) bool {
		return i < 1
	}), false)

	// Test Any
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.Any(func(i int) bool {
		return i == 1
	}), true)
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.Any(func(i int) bool {
		return i == 10
	}), false)
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.AnyParallel(func(i int) bool {
		return i == 1
	}), true)
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.AnyParallel(func(i int) bool {
		return i == 10
	}), false)

	// Test Position
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.Position(func(i int) bool {
		return i == 1
	}), optionext.Some(1))

	// Test Count
	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.Count(), 3)

	iter3 = WrapSlice(makeSlice()).Iter()
	assertEqual(t, iter3.CountParallel(), 3)

	// Test ForEach
	var j int
	WrapSlice(makeSlice()).Iter().ForEach(func(_ int) {
		j++
	})
	assertEqual(t, j, 3)

	var k int64
	WrapSlice(makeSlice()).Iter().ForEachParallel(func(_ int) {
		atomic.AddInt64(&k, 1)
	})
	assertEqual(t, k, int64(3))

	// Test Chain
	iter4 := WrapSlice(makeSlice()).Iter().Chain(WrapSlice(makeSlice()))
	assertEqual(t, iter4.Next(), optionext.Some(0))
	assertEqual(t, iter4.Next(), optionext.Some(1))
	assertEqual(t, iter4.Next(), optionext.Some(2))
	assertEqual(t, iter4.Next(), optionext.Some(0))
	assertEqual(t, iter4.Next(), optionext.Some(1))
	assertEqual(t, iter4.Next(), optionext.Some(2))
	assertEqual(t, iter4.Next(), optionext.None[int]())

	// Test Peekable
	iter5 := WrapSlice(makeSlice()).Iter().Peekable()
	assertEqual(t, iter5.Peek(), optionext.Some(0))
	assertEqual(t, iter5.Next(), optionext.Some(0))
	assertEqual(t, iter5.Peek(), optionext.Some(1))
	assertEqual(t, iter5.Next(), optionext.Some(1))
	assertEqual(t, iter5.Peek(), optionext.Some(2))
	assertEqual(t, iter5.Next(), optionext.Some(2))
	assertEqual(t, iter5.Peek(), optionext.None[int]())
	assertEqual(t, iter5.Next(), optionext.None[int]())

	// Test Reduce
	num := WrapSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}).Iter().Reduce(func(accum int, current int) int {
		return accum + current
	})
	assertEqual(t, num, optionext.Some(45))

	// Test Partition
	left, right := WrapSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}).Iter().PartitionIter(func(v int) bool {
		return v%2 == 0
	})
	assertEqual(t, left.Next(), optionext.Some(2))
	assertEqual(t, left.Next(), optionext.Some(4))
	assertEqual(t, left.Next(), optionext.Some(6))
	assertEqual(t, left.Next(), optionext.Some(8))
	assertEqual(t, left.Next(), optionext.None[int]())
	assertEqual(t, right.Next(), optionext.Some(1))
	assertEqual(t, right.Next(), optionext.Some(3))
	assertEqual(t, right.Next(), optionext.Some(5))
	assertEqual(t, right.Next(), optionext.Some(7))
	assertEqual(t, right.Next(), optionext.Some(9))
	assertEqual(t, right.Next(), optionext.None[int]())
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
