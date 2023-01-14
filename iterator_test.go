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

	// Test Map Parallel
	iter = WrapSliceMap[int, string](makeSlice()).IterPar().Map(func(v int) string {
		return strconv.Itoa(v)
	}).IterPar().CollectIter()
	assertEqual(t, iter.Len(), 3)

	// Test Filter
	iter2 := WrapSlice(makeSlice()).Iter().Filter(func(v int) bool {
		return v != 2
	}).CollectIter()
	assertEqual(t, iter2.Len(), 1)
	assertEqual(t, iter2.Next(), optionext.Some(2))
	assertEqual(t, iter2.Next(), optionext.None[int]())

	// Test Map Parallel
	iter2 = WrapSliceMap[int, string](makeSlice()).IterPar().Filter(func(v int) bool {
		return v != 2
	}).CollectIter()
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
	iter3 = WrapSlice(makeSlice()).IterPar()
	assertEqual(t, iter3.All(func(i int) bool {
		return i < 10
	}), true)
	iter3 = WrapSlice(makeSlice()).IterPar()
	assertEqual(t, iter3.All(func(i int) bool {
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
	iter3 = WrapSlice(makeSlice()).IterPar()
	assertEqual(t, iter3.Any(func(i int) bool {
		return i == 1
	}), true)
	iter3 = WrapSlice(makeSlice()).IterPar()
	assertEqual(t, iter3.Any(func(i int) bool {
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

	iter3 = WrapSlice(makeSlice()).IterPar()
	assertEqual(t, iter3.Count(), 3)

	// Test ForEach
	var j int
	WrapSlice(makeSlice()).Iter().ForEach(func(_ int) {
		j++
	})
	assertEqual(t, j, 3)

	var k int64
	WrapSlice(makeSlice()).IterPar().ForEach(func(_ int) {
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
}

func makeSlice() []int {
	return []int{0, 1, 2}
}
