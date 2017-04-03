package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoFmt(t *testing.T) {
	scenarios := []struct {
		desc        string
		in          string
		expected    string
		expectedErr bool
	}{
		{
			desc: "happy path",
			in: `package mypackage
			func main() {
				// TODO: something
			}
			`,
			expected: `package mypackage

func main() {
	// TODO: something
}
`,
			expectedErr: false,
		},
		{
			desc: "garbage in",
			in: `
			func main() {
				bad code here
			}
			`,
			expected: `
			func main() {
				bad code here
			}
			`,
			expectedErr: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := GoFmt([]byte(scenario.in))
			assert.Equal(t, scenario.expected, string(result), scenario.desc)
			assert.Equal(t, scenario.expectedErr, resultErr != nil, scenario.desc)
		})
	}

}

func TestGoImports(t *testing.T) {
	scenarios := []struct {
		desc        string
		in          string
		expected    string
		expectedErr bool
	}{
		{
			desc: "happy path - add import",
			in: `package mypackage

			func main() {
				fmt.Printf("Hello World")
			}
			`,
			expected: `package mypackage

import "fmt"

func main() {
	fmt.Printf("Hello World")
}
`,
			expectedErr: false,
		},
		{
			desc: "happy path - remove import",
			in: `package mypackage

			import (
				"fmt"
				"strings"
				)

			func main() {
				fmt.Printf("Hello World")
			}
			`,
			expected: `package mypackage

import (
	"fmt"
)

func main() {
	fmt.Printf("Hello World")
}
`,
			expectedErr: false,
		},
		{
			desc: "garbage in",
			in: `
			func main() {
				bad code here
			}
			`,
			expected: `
			func main() {
				bad code here
			}
			`,
			expectedErr: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := GoImports("myfile.go", []byte(scenario.in))
			assert.Equal(t, scenario.expected, string(result), scenario.desc)
			assert.Equal(t, scenario.expectedErr, resultErr != nil, scenario.desc)
		})
	}

}
