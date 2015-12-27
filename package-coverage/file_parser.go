package main

import (
	"log"
	"regexp"
	"strings"
)

var lineFormatChecker = regexp.MustCompile(`(?i)^(([a-z0-9.\_-]+\/)+)([a-z0-9.\_-]+).go:(([0-9]+).([0-9]+)),(([0-9]+).([0-9]+))\s[0-9]+\s[0-9]+$`)

type coverage struct {
	selfStatements int
	selfCovered    int

	childStatements int
	childCovered    int
}

func parseLines(raw string) map[string]*coverage {
	output := make(map[string]*coverage)

	lines := strings.Split(raw, "\n")
	if len(lines) == 0 {
		log.Print("no lines found in the supplied file")
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

	select {
	case <-doneCh:
		// all fragments processed
	}

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
