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

package fiximports

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyVisitor_walk(t *testing.T) {
	scenarios := []struct {
		desc             string
		input            string
		expectedStartPos int
		expectedEndPos   int
	}{
		{
			desc:             "file file 1",
			input:            testFile1,
			expectedStartPos: 14,
			expectedEndPos:   95,
		},
		{
			desc:             "file file 2",
			input:            testFile2,
			expectedStartPos: 14,
			expectedEndPos:   79,
		},
		{
			desc:             "comment at the end",
			input:            testFileCommentedImportAtEnd,
			expectedStartPos: 14,
			expectedEndPos:   75,
		},
		{
			desc:             "single import statements",
			input:            testFileIndividualImports,
			expectedStartPos: 14,
			expectedEndPos:   160,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			filename := "test.go"
			fileSet := token.NewFileSet()

			file, err := parser.ParseFile(fileSet, filename, scenario.input, parser.ParseComments)
			assert.Nil(t, err)

			visitor := newVisitor(filename, fileSet, []byte(scenario.input))
			ast.Walk(visitor, file)

			assert.Equal(t, scenario.expectedStartPos, visitor.startPos)
			assert.Equal(t, scenario.expectedEndPos, visitor.endPos)
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
			expected: `import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/corsc/go-tools/commons"
)
`,
		},
		{
			desc:  "test file 2",
			input: testFile2,
			expected: `import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
)
`,
		},
		{
			desc:  "test file - individual imports",
			input: testFileIndividualImports,
			expected: `import (
	"fmt"
	"io"
	"math"

	_ "github.com/gogo/protobuf/gogoproto"
	"github.com/golang/protobuf/proto"
)
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
				fileSet:  token.NewFileSet(),
				source:   []byte(scenario.inputFile),
				startPos: scenario.startPos,
				endPos:   scenario.endPos,
			}

			result := visitor.replaceImports(scenario.sortedImports)
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
			desc:        "test file 1",
			source:      testFile1,
			expected:    testFile1Fixed,
			expectedErr: false,
		},
		{
			desc:        "test file 2",
			source:      testFile2,
			expected:    testFile2Fixed,
			expectedErr: false,
		},
		{
			desc:        "test file 3",
			source:      testFile3,
			expected:    testFile3Fixed,
			expectedErr: false,
		},
		{
			desc:        "no imports",
			source:      testFileNoImports,
			expected:    testFileNoImports,
			expectedErr: false,
		},
		{
			desc:        "dot import",
			source:      testFileDotImport,
			expected:    testFileDotImport,
			expectedErr: false,
		},
		{
			desc:        "blank",
			source:      testFileBlankImport,
			expected:    testFileBlankImport,
			expectedErr: false,
		},
		{
			desc:        "commented above",
			source:      testFileCommentedImportAbove,
			expected:    testFileCommentedImportAbove,
			expectedErr: false,
		},
		{
			desc:        "commented at end",
			source:      testFileCommentedImportAtEnd,
			expected:    testFileCommentedImportAtEnd,
			expectedErr: false,
		},
		{
			desc:        "individual imports",
			source:      testFileIndividualImports,
			expected:    testFileIndividualImportsFixed,
			expectedErr: false,
		},
		{
			desc:        "single import",
			source:      testFileSingleImport,
			expected:    testFileSingleImportFixed,
			expectedErr: false,
		},
		{
			desc:        "extra line",
			source:      testFileExtraLine,
			expected:    testFileExtraLineFixed,
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
