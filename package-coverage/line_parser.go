package main

import (
	"fmt"
	"strconv"
	"strings"
)

type fragment struct {
	pkg        string
	statements int
	covered    bool
}

func parseLine(raw string) fragment {
	output := fragment{
		pkg: extractPackage(raw),
	}

	output.statements, output.covered = extractNumbers(raw)

	return output
}

func extractPackage(raw string) string {
	lastSlash := strings.LastIndex(raw, "/")
	if lastSlash == -1 {
		panic(fmt.Errorf("line skipped due to lack of package '%s'", raw))
	}

	return raw[:(lastSlash + 1)]
}

func extractNumbers(raw string) (int, bool) {
	parts := strings.Split(raw, " ")
	if len(parts) != 3 {
		panic(fmt.Errorf("invalid line format. parts found %d, expected 3", len(parts)))
	}

	lines := extractStatements(parts[1])
	covered := extractCovered(parts[2])

	return lines, covered
}

func extractStatements(raw string) int {
	statements, err := strconv.Atoi(raw)
	if err != nil {
		panic(err)
	}
	return statements
}

func extractCovered(raw string) bool {
	return raw == "1"
}
