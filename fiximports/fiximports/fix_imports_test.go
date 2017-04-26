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
			expectedStartPos: 23,
			expectedEndPos:   93,
		},
		{
			desc:             "file file 2",
			input:            testFile2,
			expectedStartPos: 23,
			expectedEndPos:   77,
		},
		{
			desc:             "comment at the end",
			input:            testFileCommentedImportAtEnd,
			expectedStartPos: 23,
			expectedEndPos:   73,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			visitor := &myVisitor{
				source: []byte(scenario.input),
			}

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
			expected: `	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/corsc/go-tools/commons"
`,
		},
		{
			desc:  "test file 2",
			input: testFile2,
			expected: `	"net/http"
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
				source:  []byte(scenario.inputFile),
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
		source      string
		expected    string
		expectedErr bool
	}{
		{
			desc:        "happy path - test file 1",
			source:      testFile1,
			expected:    testFile1Fixed,
			expectedErr: false,
		},
		{
			desc:        "happy path - test file 2",
			source:      testFile2,
			expected:    testFile2Fixed,
			expectedErr: false,
		},
		{
			desc:        "happy path - test file 3",
			source:      testFile3,
			expected:    testFile3Fixed,
			expectedErr: false,
		},
		{
			desc:        "happy path - no imports",
			source:      testFileNoImports,
			expected:    testFileNoImports,
			expectedErr: false,
		},
		{
			desc:        "happy path - dot import",
			source:      testFileDotImport,
			expected:    testFileDotImport,
			expectedErr: false,
		},
		{
			desc:        "happy path - blank",
			source:      testFileBlankImport,
			expected:    testFileBlankImport,
			expectedErr: false,
		},
		{
			desc:        "happy path - commented above",
			source:      testFileCommentedImportAbove,
			expected:    testFileCommentedImportAbove,
			expectedErr: false,
		},
		{
			desc:        "happy path - commented at end",
			source:      testFileCommentedImportAtEnd,
			expected:    testFileCommentedImportAtEnd,
			expectedErr: false,
		},
		// TODO: fix me
		//{
		//	desc:        "happy path - individual imports",
		//	source:      testFileIndividualImports,
		//	expected:    testFileIndividualImportsFixed,
		//	expectedErr: false,
		//},
		{
			desc:        "happy path - single import",
			source:      testFileSingleImport,
			expected:    testFileSingleImport,
			expectedErr: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := processFile("filename.go", []byte(scenario.source))

			assert.Equal(t, scenario.expectedErr, resultErr != nil, "%s %v", scenario.desc, resultErr)
			assert.Equal(t, scenario.expected, string(result), scenario.desc)
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

var testFile3 = `package main

import (
	"crypto/hmac"
	"crypto/subtle"
	"fmt"
	"math"
	"time"
	"crypto/sha256"
	"encoding/base64"
)

func main() {}
`

var testFile3Fixed = `package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"time"
)

func main() {}
`

var testFileNoImports = `package main

func main() {}
`

var testFileDotImport = `package main

import (
	"fmt"
	. "net/http/httputil"
)

func main() {}
`

var testFileBlankImport = `package main

import (
	"fmt"
	_ "net/http/httputil"
)

func main() {}
`

var testFileCommentedImportAbove = `package main

import (
	"fmt"
	// comment above
	"net/http/httputil"
)

func main() {}
`

var testFileCommentedImportAtEnd = `package main

import (
	"fmt"
	"net/http/httputil" // comment on the end
)

func main() {}
`

var testFileIndividualImports = `package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import io "io"

func main() {}
`

var testFileIndividualImportsFixed = `package main

import (
	"io"
	"fmt"
	"math"

	_ "github.com/gogo/protobuf/gogoproto"
	"github.com/golang/protobuf/proto"
)

func main() {}
`

var testFileSingleImport = `package main

import "github.com/golang/protobuf/proto"

func main() {}
`
