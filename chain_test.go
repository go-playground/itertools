package itertools

import (
	optionext "github.com/go-playground/pkg/v5/values/option"
	"testing"
)

func TestChain(t *testing.T) {
	// Test sort
	slice := []int{0, 1, 2, 3}
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

	// Test same Iterator[T] but different underlying iterator types
	fi := &fakeIterator{max: 3}
	si := WrapSlice(slice).IntoIter()

	iter := Chain[int](fi, si)
	assertEqual(t, iter.Next(), optionext.Some(2))
	assertEqual(t, iter.Next(), optionext.Some(1))
	assertEqual(t, iter.Next(), optionext.Some(0))
	assertEqual(t, iter.Next(), optionext.Some(0))
	assertEqual(t, iter.Next(), optionext.Some(1))
	assertEqual(t, iter.Next(), optionext.Some(2))
	assertEqual(t, iter.Next(), optionext.Some(3))
	assertEqual(t, iter.Next(), optionext.None[int]())
}

type fakeIterator struct {
	max int
}

func (f *fakeIterator) Next() optionext.Option[int] {
	f.max--
	if f.max < 0 {
		return optionext.None[int]()
	}
	return optionext.Some(f.max)
}
