package main

import (
	"fmt"
	"github.com/go-playground/itertools"
)

func main() {
	results := itertools.WrapSlice([]int{4, 3, 2, 1, 0}).Iter().Filter(func(v int) bool {
		if v >= 5 {
			return true
		}
		return false
	}).Collect()

	fmt.Printf("%#v\n", results)
}
