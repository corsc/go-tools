package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateCoverage(t *testing.T) {
	filename := "../test-data/acc.out"

	expected := map[string]*coverage{
		"sage42.org/go-tools/package-coverage/": {
			selfStatements: 87,
			selfCovered:    55,

			childStatements: 0,
			childCovered:    0,
		},
	}

	result := CalculateCoverage(filename)
	converted := map[string]*coverage(result)
	assert.Equal(t, expected, converted)
}
