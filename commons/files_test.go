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

package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected bool
	}{
		{
			desc:     "exists",
			input:    "./testdata/file-exists/exists.go",
			expected: true,
		},
		{
			desc:     "doesn't exist",
			input:    "./testdata/file-exists/doesnt-exist.go",
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result := FileExists(scenario.input)
			assert.Equal(t, scenario.expected, result, scenario.desc)
		})
	}

}

func TestGetGoFiles(t *testing.T) {
	scenarios := []struct {
		desc        string
		in          []string
		expected    []string
		expectedErr bool
	}{
		{
			desc:        "empty input",
			in:          []string{},
			expected:    []string{},
			expectedErr: false,
		},
		{
			desc:        "happy path",
			in:          []string{"apples.go"},
			expected:    []string{"apples.go"},
			expectedErr: false,
		},
		{
			desc:        "bad file",
			in:          []string{"apples.txt"},
			expected:    nil,
			expectedErr: true,
		},
		{
			desc:        "happy path - many",
			in:          []string{"apples.go", "oranges.go"},
			expected:    []string{"apples.go", "oranges.go"},
			expectedErr: false,
		},
		{
			desc:        "bad file",
			in:          []string{"apples.go", "orangos.txt"},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := GetGoFiles(scenario.in...)
			assert.Equal(t, scenario.expected, result, scenario.desc)
			assert.Equal(t, scenario.expectedErr, (resultErr != nil), scenario.desc)
		})
	}
}

func TestGetGoFilesFromDir(t *testing.T) {
	scenarios := []struct {
		desc      string
		input     string
		expected  []string
		expectErr bool
	}{
		{
			desc:  "valid directory",
			input: "./testdata/get-go-files/",
			expected: []string{
				"testdata/get-go-files/a.go",
				"testdata/get-go-files/b.go",
				"testdata/get-go-files/c_test.go",
			},
			expectErr: false,
		},
		{
			desc:      "invalid directory",
			input:     "/something/invalid",
			expected:  []string{},
			expectErr: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := GetGoFilesFromDir(scenario.input)
			assert.Equal(t, scenario.expected, result)
			assert.Equal(t, scenario.expectErr, resultErr != nil)
		})
	}
}

func TestGetGoFilesFromDirectoryRecursive(t *testing.T) {
	scenarios := []struct {
		desc      string
		input     string
		expected  []string
		expectErr bool
	}{
		{
			desc:  "valid dir with ...",
			input: "./testdata/get-go-files/...",
			expected: []string{
				"testdata/get-go-files/a.go",
				"testdata/get-go-files/b.go",
				"testdata/get-go-files/c_test.go",
				"testdata/get-go-files/c/d.go",
				"testdata/get-go-files/c/e.go",
			},
			expectErr: false,
		},
		{
			desc:  "valid dir without ...",
			input: "./testdata/get-go-files/",
			expected: []string{
				"testdata/get-go-files/a.go",
				"testdata/get-go-files/b.go",
				"testdata/get-go-files/c_test.go",
				"testdata/get-go-files/c/d.go",
				"testdata/get-go-files/c/e.go",
			},
			expectErr: false,
		},
		{
			desc:      "invalid dir with ...",
			input:     "./something/invalid/",
			expected:  []string{},
			expectErr: true,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := GetGoFilesFromDirectoryRecursive(scenario.input)
			assert.Equal(t, scenario.expected, result, scenario.desc)
			assert.Equal(t, scenario.expectErr, resultErr != nil, scenario.desc)
		})
	}
}
