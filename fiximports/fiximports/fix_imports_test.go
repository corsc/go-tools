package fiximports

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyVisitor_getImportBoundaries(t *testing.T) {
	scenarios := []struct {
		desc             string
		input            string
		expectedStartPos int
		expectedEndPos   int
	}{
		{
			desc:             "file file 1",
			input:            testFile1,
			expectedStartPos: 24,
			expectedEndPos:   93,
		},
		{
			desc:             "file file 2",
			input:            testFile2,
			expectedStartPos: 24,
			expectedEndPos:   77,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			visitor := &myVisitor{}

			file, err := parser.ParseFile(token.NewFileSet(), "test.go", scenario.input, parser.ParseComments)
			assert.Nil(t, err)

			start, end := visitor.getImportBoundaries(file)
			assert.Equal(t, scenario.expectedStartPos, start)
			assert.Equal(t, scenario.expectedEndPos, end)
		})
	}
}

func TestMyVisitor_orderImports(t *testing.T) {
	visitor := &myVisitor{}

	file, err := parser.ParseFile(token.NewFileSet(), "test.go", testFile1, parser.ParseComments)
	assert.Nil(t, err)

	visitor.orderImports(file)

	expected := []string{
		`"flag"`,
		`"fmt"`,
		`"github.com/corsc/go-tools/commons"`,
		`"os"`,
		`"strings"`,
	}

	for index, path := range expected {
		assert.Equal(t, path, file.Imports[index].Path.Value)
	}
}

func TestMyVisitor_generateImportsFragment(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:  "test file 1",
			input: testFile1,
			expected: `"flag"
	"fmt"
	"os"
	"strings"

	"github.com/corsc/go-tools/commons"
`,
		},
		{
			desc:  "test file 2",
			input: testFile2,
			expected: `"net/http"
	"net/http/httptest"
	"net/http/httputil"
`,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			visitor := &myVisitor{}

			file, err := parser.ParseFile(token.NewFileSet(), "test.go", scenario.input, parser.ParseComments)
			assert.Nil(t, err)

			visitor.orderImports(file)
			result := visitor.generateImportsFragment(file)

			assert.Equal(t, scenario.expected, result, scenario.desc)
		})
	}
}

func TestMyVisitor_replaceImports(t *testing.T) {
	scenarios := []struct {
		desc          string
		inputFile     string
		sortedImports string
		startPos      int
		endPos        int
		expected      string
	}{
		{
			desc:      "test file 1",
			inputFile: testFile1,
			sortedImports: `"flag"
	"fmt"
	"os"
	"strings"

	"github.com/corsc/go-tools/commons"
`,
			startPos: 24,
			endPos:   93,
			expected: testFile1Fixed,
		},
		{
			desc:      "test file 2",
			inputFile: testFile2,
			sortedImports: `"net/http"
	"net/http/httptest"
	"net/http/httputil"
`,
			startPos: 24,
			endPos:   77,
			expected: testFile2Fixed,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			visitor := &myVisitor{
				fileSet: token.NewFileSet(),
			}

			file, err := parser.ParseFile(visitor.fileSet, "test.go", scenario.inputFile, parser.ParseComments)
			assert.Nil(t, err)

			result := visitor.replaceImports(file, scenario.sortedImports, scenario.startPos, scenario.endPos)

			assert.Equal(t, scenario.expected, string(result), scenario.desc)
		})
	}
}

func TestProcessFile(t *testing.T) {
	scenarios := []struct {
		desc        string
		source      []byte
		expected    []byte
		expectedErr bool
	}{
		{
			desc:        "happy path - test file 1",
			source:      []byte(testFile1),
			expected:    []byte(testFile1Fixed),
			expectedErr: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := processFile("filename.go", scenario.source)
			assert.Equal(t, scenario.expected, result, scenario.desc)
			assert.Equal(t, scenario.expectedErr, resultErr != nil, scenario.desc)
		})
	}

}

var testFile1 = `package main

import (
	"github.com/corsc/go-tools/commons"

	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {}
`

var testFile1Fixed = `package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/corsc/go-tools/commons"
)

func main() {}
`

var testFile2 = `package main

import (
	"net/http/httptest"
	"net/http"
	"net/http/httputil"
)

func main() {}
`

var testFile2Fixed = `package main

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
)

func main() {}
`
