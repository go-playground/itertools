package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"testing"
)

func TestMap(t *testing.T) {
	// Test Next
	iter := MapIter(makeMap())
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), true)
	assertEqual(t, iter.Next().IsSome(), false)

	// Test Retain
	iter = MapIter(makeMap()).Retain(func(k string, v int) bool {
		return v == 3
	})
	assertEqual(t, iter.Next(), optionext.Some(Entry[string, int]{Key: "3", Value: 3}))
	assertEqual(t, iter.Next(), optionext.None[Entry[string, int]]())
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
