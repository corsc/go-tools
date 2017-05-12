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
	"fmt"
	"regexp"
	"strings"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// coverage data line format (used to filter out the rest of the coverage file contents)
const lineFormat = `(?i)^(([a-z0-9.\_-]+\/)+)([a-z0-9.\_-]+).go:(([0-9]+).([0-9]+)),(([0-9]+).([0-9]+))\s[0-9]+\s[0-9]+$`

var lineFormatChecker = regexp.MustCompile(lineFormat)

type coverage struct {
	selfStatements int
	selfCovered    int

	childStatements int
	childCovered    int
}

func (c *coverage) String() string {
	return fmt.Sprintf("self: %d/%d / child: %d/%d", c.selfCovered, c.selfStatements, c.childCovered, c.childStatements)
}

// convert string contents of the coverage files into data structures
func parseLines(raw string) map[string]*coverage {
	output := make(map[string]*coverage)

	lines := strings.Split(raw, "\n")
	if len(lines) == 0 {
		utils.LogWhenVerbose("[print] no lines found in the supplied file")
		return output
	}

	fragmentCh := make(chan fragment, len(lines))
	doneCh := processFragments(output, fragmentCh)

	for _, line := range lines {
		if !validLineFormat(line) {
			continue
		}

		fragmentCh <- parseLine(line)
	}
	close(fragmentCh)

	<-doneCh

	return output
}

func validLineFormat(line string) bool {
	return lineFormatChecker.MatchString(line)
}

func processFragments(output map[string]*coverage, fragmentCh chan fragment) chan struct{} {
	doneCh := make(chan struct{})

	go func() {
		for fragment := range fragmentCh {
			coverage := getOrCreateCoverage(output, fragment.pkg)
			processSelfCoverage(coverage, fragment)
		}

		updateChildCoverage(output)

		close(doneCh)
	}()

	return doneCh
}

func getOrCreateCoverage(output map[string]*coverage, pkg string) *coverage {
	cover, ok := output[pkg]
	if !ok {
		output[pkg] = &coverage{}
		cover = output[pkg]
	}
	return cover
}

func processSelfCoverage(cover *coverage, fragment fragment) {
	cover.selfStatements += fragment.statements
	if fragment.covered {
		cover.selfCovered += fragment.statements
	}
}

// TODO: make this more efficient
func updateChildCoverage(output map[string]*coverage) {
	for pkgOuter, coverageOuter := range output {
		for pkgInner, coverageInner := range output {
			if pkgOuter == pkgInner {
				continue
			}

			if isChild(pkgOuter, pkgInner) {
				coverageOuter.childStatements += coverageInner.selfStatements
				coverageOuter.childCovered += coverageInner.selfCovered
			}
		}
	}
}

func isChild(pkgA string, pkgB string) bool {
	return strings.HasPrefix(pkgB, pkgA)
}
