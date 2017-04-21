package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllPackagesUnderDirectory(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected []string
	}{
		{
			desc:  "happy path",
			input: "./testdata/get-go-files/",
			expected: []string{
				"./testdata/get-go-files",
			},
		},
		{
			desc:  "happy path - recursive",
			input: "./testdata/get-go-files/...",
			expected: []string{
				"./testdata/get-go-files",
				"./testdata/get-go-files/c",
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result := GetAllPackagesUnderDirectory(scenario.input)
			assert.Equal(t, scenario.expected, result, scenario.desc)
		})
	}
}
