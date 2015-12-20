package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	in := "sage42.org/go-tools/package-coverage/line_parser.go:9.37,11.2 1 0"
	expected := fragment{
		pkg:     "sage42.org/go-tools/package-coverage/",
		lines:   1,
		covered: 0,
	}

	result := parseLine(in)
	assert.Equal(t, expected, result)
}

func TestExtractPackage_HappyPath(t *testing.T) {
	in := "sage42.org/go-tools/package-coverage/line_parser.go:9.37,11.2 1 0"
	expected := "sage42.org/go-tools/package-coverage/"

	result := extractPackage(in)
	assert.Equal(t, expected, result)
}

func TestExtractPackage_InvalidLinePanic(t *testing.T) {
	in := ""

	assert.Panics(t, func() {
		extractPackage(in)
	})
}

func TestExtractLines_HappyPath(t *testing.T) {
	in := "sage42.org/go-tools/package-coverage/line_parser.go:9.37,11.2 1 0"
	expectedLines := 1
	expectedCovered := 0

	resultLines, resultCoverted := extractLines(in)
	assert.Equal(t, expectedLines, resultLines)
	assert.Equal(t, expectedCovered, resultCoverted)
}

func TestExtractLines_InvalidLinePanics(t *testing.T) {
	in := ""

	assert.Panics(t, func() {
		extractLines(in)
	})
}
