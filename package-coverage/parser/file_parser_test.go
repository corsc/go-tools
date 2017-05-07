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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLines_OnePackage(t *testing.T) {
	in := `mode: set
github.com/corsc/go-tools/package-coverage/file_parser.go:13.49,15.2 1 0
github.com/corsc/go-tools/package-coverage/line_parser.go:15.37,23.2 3 1
github.com/corsc/go-tools/package-coverage/line_parser.go:25.40,27.21 2 1
`
	expected := map[string]*coverage{
		"github.com/corsc/go-tools/package-coverage/": {
			selfStatements: 6,
			selfCovered:    5,
		},
	}

	result := parseLines(in)
	assert.Equal(t, expected, result)
}

func TestParseLines_TwoPackages(t *testing.T) {
	in := `mode: set
github.com/corsc/go-tools/package-coverage/file_parser.go:13.49,15.2 1 0
github.com/corsc/go-tools/some-other-package/something.go:25.40,27.21 2 1
`
	expected := map[string]*coverage{
		"github.com/corsc/go-tools/package-coverage/": {
			selfStatements: 1,
			selfCovered:    0,
		},
		"github.com/corsc/go-tools/some-other-package/": {
			selfStatements: 2,
			selfCovered:    2,
		},
	}

	result := parseLines(in)
	assert.Equal(t, expected, result)
}

func TestParseLines_PackageAndChild(t *testing.T) {
	in := `mode: set
github.com/corsc/go-tools/package-coverage/file_parser.go:13.49,15.2 1 1
github.com/corsc/go-tools/package-coverage/sub/file_parser.go:13.49,15.2 1 1
github.com/corsc/go-tools/package-coverage/sub/other.go:13.49,15.2 1 0
`
	expected := map[string]*coverage{
		"github.com/corsc/go-tools/package-coverage/": {
			selfStatements:  1,
			selfCovered:     1,
			childStatements: 2,
			childCovered:    1,
		},
		"github.com/corsc/go-tools/package-coverage/sub/": {
			selfStatements: 2,
			selfCovered:    1,
		},
	}

	result := parseLines(in)
	assert.Equal(t, expected, result)
}

func TestValidLineFormat(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected bool
	}{
		{
			desc:     "invalid line - empty string",
			input:    "",
			expected: false,
		},
		{
			desc:     "invalid line - blank string",
			input:    "    ",
			expected: false,
		},
		{
			desc:     "invalid line - mode line",
			input:    "mode: set",
			expected: false,
		},
		{
			desc:     "invalid line - random content",
			input:    ";lfksal;fkaof'eqr'pjasnvaht;8ehtgnq;s",
			expected: false,
		},
		{
			desc:     "valid line - properly formatted line",
			input:    "github.com/corsc/go-tools/package-coverage/line_parser.go:54.38,56.2 1 1",
			expected: true,
		},
		{
			desc:     "valid line - strange case line",
			input:    "github.com/corsc/go-tools/package-coverage/line_parser.go:54.38,56.2 1 1",
			expected: true,
		},
	}

	for _, scenario := range scenarios {
		result := validLineFormat(scenario.input)

		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}

func TestIsChild(t *testing.T) {
	scenarios := []struct {
		desc     string
		inA      string
		inB      string
		expected bool
	}{
		{
			desc:     "Not a child",
			inA:      "fu/bar",
			inB:      "google.com",
			expected: false,
		},
		{
			desc:     "Is a direct child",
			inA:      "fu.bar/com/",
			inB:      "fu.bar/com/fu/",
			expected: true,
		},
		{
			desc:     "Is a sub child",
			inA:      "fu.bar/com/",
			inB:      "fu.bar/com/fu/bar/",
			expected: true,
		},
	}

	for _, scenario := range scenarios {
		result := isChild(scenario.inA, scenario.inB)
		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}
