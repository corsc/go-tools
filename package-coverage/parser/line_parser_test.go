package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	in := "github.com/corsc/go-tools/package-coverage/line_parser.go:9.37,11.2 1 0"
	expected := fragment{
		pkg:        "github.com/corsc/go-tools/package-coverage/",
		statements: 1,
		covered:    false,
	}

	result := parseLine(in)
	assert.Equal(t, expected, result)
}

func TestExtractPackage_HappyPath(t *testing.T) {
	in := "github.com/corsc/go-tools/package-coverage/line_parser.go:9.37,11.2 1 0"
	expected := "github.com/corsc/go-tools/package-coverage/"

	result := extractPackage(in)
	assert.Equal(t, expected, result)
}

func TestExtractPackage_InvalidLinePanic(t *testing.T) {
	in := ""

	assert.Panics(t, func() {
		extractPackage(in)
	})
}

func TestExtractNumbers_HappyPath(t *testing.T) {
	in := "github.com/corsc/go-tools/package-coverage/line_parser.go:9.37,11.2 1 0"
	expectedStatements := 1
	expectedCovered := false

	resultLines, resultCoverted := extractNumbers(in)
	assert.Equal(t, expectedStatements, resultLines)
	assert.Equal(t, expectedCovered, resultCoverted)
}

func TestExtractNumbers_InvalidLinePanics(t *testing.T) {
	in := ""

	assert.Panics(t, func() {
		extractNumbers(in)
	})
}

func TestExtractStatements_HappyPath(t *testing.T) {
	scenarios := []struct {
		in       string
		expected int
	}{
		{
			in:       "0",
			expected: 0,
		},
		{
			in:       "666",
			expected: 666,
		},
	}

	for _, scenario := range scenarios {
		result := extractStatements(scenario.in)
		assert.Equal(t, scenario.expected, result)
	}
}

func TestExtractStatements_InvalidInputPanics(t *testing.T) {
	in := ""

	assert.Panics(t, func() {
		extractStatements(in)
	})
}

func TestExtractCovered_HappyPath(t *testing.T) {
	scenarios := []struct {
		in       string
		expected bool
	}{
		{
			in:       "0",
			expected: false,
		},
		{
			in:       "1",
			expected: true,
		},
	}

	for _, scenario := range scenarios {
		result := extractCovered(scenario.in)
		assert.Equal(t, scenario.expected, result)
	}
}
