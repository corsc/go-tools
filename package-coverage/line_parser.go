package main

import (
	"fmt"
	"strconv"
	"strings"
)

type fragment struct {
	pkg     string
	lines   int
	covered int
}

func parseLine(raw string) fragment {
	output := fragment{
		pkg: extractPackage(raw),
	}

	output.lines, output.covered = extractLines(raw)

	return output
}

func extractPackage(raw string) string {
	lastSlash := strings.LastIndex(raw, "/")
	if lastSlash == -1 {
		panic(fmt.Errorf("line skipped due to lack of package '%s'", raw))
	}

	return raw[:(lastSlash + 1)]
}

func extractLines(raw string) (int, int) {
	parts := strings.Split(raw, " ")
	if len(parts) != 3 {
		panic(fmt.Errorf("invalid line format. parts found %d, expected 3"))
	}

	lines, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	covered, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}

	return lines, covered
}
