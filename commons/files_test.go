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
	input := "./testdata/get-go-files/"
	expected := []string{
		"testdata/get-go-files/a.go",
		"testdata/get-go-files/b.go",
	}

	result, resultErr := GetGoFilesFromDir(input)
	assert.Equal(t, expected, result)
	assert.Nil(t, resultErr)
}

func TestGetGoFilesFromDirectoryRecursive(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected []string
	}{
		{
			desc:  "valid dir with ...",
			input: "./testdata/get-go-files/...",
			expected: []string{
				"testdata/get-go-files/a.go",
				"testdata/get-go-files/b.go",
				"testdata/get-go-files/c/d.go",
				"testdata/get-go-files/c/e.go",
			},
		},
		{
			desc:  "valid dir without ...",
			input: "./testdata/get-go-files/",
			expected: []string{
				"testdata/get-go-files/a.go",
				"testdata/get-go-files/b.go",
				"testdata/get-go-files/c/d.go",
				"testdata/get-go-files/c/e.go",
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := GetGoFilesFromDirectoryRecursive(scenario.input)
			assert.Equal(t, scenario.expected, result, scenario.desc)
			assert.Nil(t, resultErr, scenario.desc)
		})
	}
}
