package itertools

import (
	. "github.com/go-playground/assert/v2"
	optionext "github.com/go-playground/pkg/v5/values/option"
	"testing"
)

func TestChain(t *testing.T) {
	// Test sort
	slice := []int{0, 1, 2, 3}
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

	iterChain = Chain[int](WrapSlice(slice).IntoIter(), WrapSlice(slice).IntoIter()).Iter()
	Equal(t, iterChain.Next(), optionext.Some(0))
	Equal(t, iterChain.Next(), optionext.Some(1))
	Equal(t, iterChain.Next(), optionext.Some(2))
	Equal(t, iterChain.Next(), optionext.Some(3))
	Equal(t, iterChain.Next(), optionext.Some(0))
	Equal(t, iterChain.Next(), optionext.Some(1))
	Equal(t, iterChain.Next(), optionext.Some(2))
	Equal(t, iterChain.Next(), optionext.Some(3))
	Equal(t, iterChain.Next(), optionext.None[int]())

	// Test same Iterator[T] but different underlying iterator types
	fi := &fakeIterator{max: 3}
	si := WrapSlice(slice).IntoIter()

	iter := Chain[int](fi, si)
	Equal(t, iter.Next(), optionext.Some(2))
	Equal(t, iter.Next(), optionext.Some(1))
	Equal(t, iter.Next(), optionext.Some(0))
	Equal(t, iter.Next(), optionext.Some(0))
	Equal(t, iter.Next(), optionext.Some(1))
	Equal(t, iter.Next(), optionext.Some(2))
	Equal(t, iter.Next(), optionext.Some(3))
	Equal(t, iter.Next(), optionext.None[int]())
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
