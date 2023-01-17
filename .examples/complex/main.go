package main

import (
	"fmt"
	"github.com/go-playground/itertools"
	optionext "github.com/go-playground/pkg/v5/values/option"
	"strconv"
)

type FakeIterator struct {
	max int
}

func (f *FakeIterator) Next() optionext.Option[int] {
	f.max--
	if f.max < 0 {
		return optionext.None[int]()
	}
	return optionext.Some(f.max)
}

func main() {
	results := itertools.WrapSliceMap[int, string]([]int{4, 3, 2, 1, 0}).Iter().Chain(&FakeIterator{
		max: 10,
	}).Filter(func(v int) bool {
		if v >= 5 {
			return true
		}
		return false
	}).StepBy(2).Take(6).Map(func(v int) string {
		return strconv.Itoa(v)
	}).Iter().Collect()

	fmt.Printf("%#v\n", results)
}
