goja-babel
==========

[![CI](https://github.com/jvatic/goja-babel/actions/workflows/ci.yml/badge.svg)](https://github.com/jvatic/goja-babel/actions/workflows/ci.yml)

Uses [goja](https://github.com/dop251/goja) to run [babel.js](https://babeljs.io/) within Go.

## Usage

```go
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jvatic/goja-babel"
)

func main() {
	babel.Init(4) // Setup 4 transformers (can be any number > 0)
	res, err := babel.Transform(strings.NewReader(`let foo = 1;
	<div>
		Hello JSX!
		The value of foo is {foo}.
	</div>`), map[string]interface{}{
		"plugins": []string{
			"transform-react-jsx",
			"transform-block-scoping",
		},
	})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, res)
	fmt.Println("")
}
```

```js
$ go run main.go
import { jsxs as _jsxs } from "react/jsx-runtime";
var foo = 1;
/*#__PURE__*/_jsxs("div", {
  children: ["Hello JSX! The value of foo is ", foo, "."]
});
```

## Benchmarks

```
go test -bench Transform -benchmem
goos: darwin
goarch: arm64
pkg: github.com/jvatic/goja-babel
cpu: Apple M3 Pro
BenchmarkTransformString-11                     	     402	   2984801 ns/op	 2310079 B/op	   29834 allocs/op
BenchmarkTransformStringWithSingletonPool-11    	     411	   2911610 ns/op	 2310858 B/op	   29835 allocs/op
BenchmarkTransformStringWithLargePool-11        	     411	   2982102 ns/op	 2310419 B/op	   29835 allocs/op
PASS
ok  	github.com/jvatic/goja-babel	3.878s
```
