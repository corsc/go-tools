// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gotools

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
			desc: "remove import",
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
