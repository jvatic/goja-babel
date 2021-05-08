package babel

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stvp/assert"
)

var opts = map[string]interface{}{
	"plugins": []string{
		"transform-react-jsx",
		"transform-block-scoping",
	},
}

func _TransformWithPool(t *testing.T, n int, p int) {
	input := `let foo = 1;
	<div>
		Hello JSX!
		The value of foo is {foo}.
	</div>`
	expectedOutput := strings.Join([]string{
		"var foo = 1;",
		"",
		"/*#__PURE__*/",
		`React.createElement("div", null, "Hello JSX! The value of foo is ", foo, ".");`,
	}, "\n")
	done := make(chan struct{})
	if p > 0 {
		Init(p)
	}
	for i := 0; i < n; i++ { // make sure pool works by calling Transform multiple times
		go func(i int) {
			defer func() { done <- struct{}{} }()
			output, err := Transform(strings.NewReader(input), opts)
			assert.Nil(t, err, fmt.Sprintf("Transform(%d) failed", i))
			var outputBuf bytes.Buffer
			io.Copy(&outputBuf, output)
			assert.Equal(t, expectedOutput, outputBuf.String(), fmt.Sprintf("Transform(%d) failed", i))
		}(i)
	}
	for i := 0; i < n; i++ {
		<-done
	}
}

func TestTransform(t *testing.T) {
	_TransformWithPool(t, 2, 0)
}

func TestTransformWithPool(t *testing.T) {
	_TransformWithPool(t, 4, 2)
}

func BenchmarkTransformString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TransformString(`let foo = 1`, opts)
	}
}

func BenchmarkTransformStringWithSingletonPool(b *testing.B) {
	Init(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TransformString(`let foo = 1`, opts)
	}
}

func BenchmarkTransformStringWithLargePool(b *testing.B) {
	Init(4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TransformString(`let foo = 1`, opts)
	}
}
