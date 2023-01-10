package main

import (
	"fmt"
	iterext "github.com/go-playground/pkg/v5/iter"
)

func main() {
	results := iterext.SliceIter([]int{4, 3, 2, 1, 0}).Iter().Filter(func(v int) bool {
		if v >= 5 {
			return true
		}
		return false
	}).Collect()

	fmt.Printf("%#v\n", results)
}
