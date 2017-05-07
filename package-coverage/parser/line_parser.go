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
	"log"
	"strconv"
	"strings"
)

type fragment struct {
	pkg        string
	file       string
	statements int
	covered    bool
}

func parseLine(raw string) fragment {
	output := fragment{
		pkg:  extractPackage(raw),
		file: extractFile(raw),
	}

	output.statements, output.covered = extractNumbers(raw)

	return output
}

func extractPackage(raw string) string {
	lastSlash := strings.LastIndex(raw, "/")
	if lastSlash == -1 {
		log.Panicf("line skipped due to lack of package '%s'", raw)
	}

	return raw[:(lastSlash + 1)]
}

func extractFile(raw string) string {
	lastSlash := strings.LastIndex(raw, "/")
	if lastSlash == -1 {
		log.Panicf("line skipped due to lack of package '%s'", raw)
	}

	fileAndLines := raw[(lastSlash + 1):]
	line := strings.LastIndex(fileAndLines, ":")
	if line == -1 {
		log.Panicf("line skipped due to lack of line number '%s'", raw)
	}

	return fileAndLines[:line]
}

func extractNumbers(raw string) (int, bool) {
	parts := strings.Split(raw, " ")
	if len(parts) != 3 {
		log.Panicf("invalid line format. parts found %d, expected 3", len(parts))
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
	covered, err := strconv.Atoi(raw)
	if err != nil {
		panic(err)
	}
	return covered > 0
}
