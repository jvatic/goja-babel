package babel

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stvp/assert"
)

var opts = map[string]any{
	"plugins": []string{
		"transform-react-jsx",
		"transform-block-scoping",
	},
}

var input = `let foo = 1;
	<div>
		Hello JSX!
		The value of foo is {foo}.
	</div>`

func _TransformWithPool(t *testing.T, n int, p int) {
	expectedOutput := strings.Join([]string{
		`import { jsxs as _jsxs } from "react/jsx-runtime";`,
		"var foo = 1;",
		`/*#__PURE__*/_jsxs("div", {`,
		`  children: ["Hello JSX! The value of foo is ", foo, "."]`,
		`});`,
	}, "\n")
	if p > 0 {
		Init(p)
	}
	type result struct {
		Index  int
		Error  error
		Output string
	}
	outputCh := make(chan result)
	for i := range n { // make sure pool works by calling Transform multiple times
		go func(i int) {
			output, err := Transform(strings.NewReader(input), opts)
			if err != nil {
				outputCh <- result{
					Index: i,
					Error: err,
				}
				return
			}
			var outputBuf bytes.Buffer
			io.Copy(&outputBuf, output)
			outputCh <- result{
				Index:  i,
				Output: outputBuf.String(),
			}
		}(i)
	}

	for range n { // check outputs
		output := <-outputCh
		assert.Nil(t, output.Error, fmt.Sprintf("Transform(%d) failed", output.Index))
		assert.Equal(t, expectedOutput, output.Output, fmt.Sprintf("Transform(%d) failed", output.Index))
	}
}

func TestTransform(t *testing.T) {
	_TransformWithPool(t, 2, 0)
}

func TestTransformWithPool(t *testing.T) {
	_TransformWithPool(t, 4, 2)
}

func BenchmarkTransformString(b *testing.B) {
	for b.Loop() {
		TransformString(input, opts)
	}
}

func BenchmarkTransformStringWithSingletonPool(b *testing.B) {
	Init(1)

	for b.Loop() {
		TransformString(input, opts)
	}
}

func BenchmarkTransformStringWithLargePool(b *testing.B) {
	Init(4)

	for b.Loop() {
		TransformString(input, opts)
	}
}
