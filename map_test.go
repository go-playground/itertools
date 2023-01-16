package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"testing"
)

func TestMap(t *testing.T) {

	// Test Misc
	iter := WrapMap(makeMap())
	assertEqual(t, iter.Len(), 5)

	// Test Next
	iter = WrapMap(makeMap())
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), false)

	// Test Retain
	iter = WrapMap(makeMap()).Retain(func(entry Entry[string, int]) bool {
		return entry.Value == 3
	})
	assertEqual(t, iter.Next(), optionext.Some(Entry[string, int]{Key: "3", Value: 3}))
	assertEqual(t, iter.Next(), optionext.None[Entry[string, int]]())

	// Test Retain Function
	m := makeMap()
	RetainMap(m, func(entry Entry[string, int]) bool {
		return entry.Value == 3
	})
	assertEqual(t, len(m), 1)
	assertEqual(t, m["3"], 3)

	// Test Iter Filter
	iter2 := WrapMap(makeMap()).Iter().Filter(func(v Entry[string, int]) bool {
		return v.Value != 3
	})
	assertEqual(t, iter2.Next(), optionext.Some(Entry[string, int]{Key: "3", Value: 3}))
	assertEqual(t, iter2.Next(), optionext.None[Entry[string, int]]())
}

func makeMap() map[string]int {
	return map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
	}
}

func BenchmarkRetainMap_Retain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RetainMap(makeMap(), func(entry Entry[string, int]) (retain bool) {
			return entry.Value == 3
		})
	}
}
