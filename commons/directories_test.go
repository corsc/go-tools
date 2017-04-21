package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected bool
	}{
		{
			desc:     "happy path",
			input:    "./testdata/file-exists/",
			expected: true,
		},
		{
			desc:     "doesn't exist",
			input:    "./testdata/doesnt-exist/",
			expected: false,
		},
		{
			desc:     "not a directory",
			input:    "./testdata/file-exists/exists.go",
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result := IsDir(scenario.input)
			assert.Equal(t, scenario.expected, result, scenario.desc)
		})
	}

}
