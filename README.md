# itertools

![Project status](https://img.shields.io/badge/version-0.1.0-green.svg)
[![GoDoc](https://godoc.org/github.com/go-playground/itertools?status.svg)](https://pkg.go.dev/mod/github.com/go-playground/itertools)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Go Iteration tools with a rusty flavour

## Requirements

- Go 1.19+

## Motivation

1. I missed some iteration style tools in Rust.
2. Wanted to see how far I could push Go generics(turns out it very limited).
3. For fun.

## Example Usage (Simple)

```go
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
```

## Example Usage (Complex)

See examples for more complex usages.

```go
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
```

## Caveats

- `Map` and it's `MAP` type parameter must be defined and passed around to be able to using inline. This is a limitation of Go generics not allowing new type parameters on methods.
- `Chunk` can only be used at the end of a series of iterators from `Iter` but can be used and wrapped by `Iter` again. This is a limitation of the Go Compiler which causes a recursive initialization issue https://github.com/golang/go/issues/50215.
- `Iter` must be called on some types, like the wrapped slice or map types, to allow usage of helper functions tied directly to them but not `Iterate`
- `Reduce` will have to be used directly instead of other helper function such as `Sum`, `Product`, ... again because no new type parameters on methods.

## How to Contribute

Make a pull request... can't guarantee it will be added, going to strictly vet what goes in.

## License

Distributed under MIT License, please see license file within the code for more details.
