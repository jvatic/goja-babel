goja-babel
==========

Uses github.com/dop251/goja to run babel.js within Go.

**WARNING:** This is largely untested and the exposed API may change at any time.

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
	res, err := babel.Transform(strings.NewReader(`let foo = 1;
	<div>
		Hello JSX!
		The value of foo is {foo}.
	</div>`), map[string]interface{}{
		"plugins": []string{
			"transform-react-jsx",
			"transform-es2015-block-scoping",
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
var foo = 1;
React.createElement(
	"div",
	null,
	"Hello JSX! The value of foo is ",
	foo,
	"."
);
```
