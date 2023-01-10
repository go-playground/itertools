# itertools

![Project status](https://img.shields.io/badge/version-0.1.0-green.svg)
[![GoDoc](https://godoc.org/github.com/go-playground/itertools?status.svg)](https://pkg.go.dev/mod/github.com/go-playground/itertools)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Go Iteration tools with a rusty flavour

## Motivation

1. I missed some of the iteration style tools in Rust.
2. Wanted to see how far I could push Go generics.
3. For fun.

## Example Usage (Simple)

See examples for more complex usages.

```go
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

```

## How to Contribute

Make a pull request... can't guarantee it will be added, going to strictly vet what goes in.

## License

Distributed under MIT License, please see license file within the code for more details.
