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
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// CoverageByPackage contains the result of parsing one or more package's coverage file
type coverageByPackage map[string]*coverage

// PrintCoverage will attempt to calculate and print the coverage from the supplied coverage file to standard out.
func PrintCoverage(writer io.Writer, basePath string, exclusionsMatcher *regexp.Regexp, minCoverage int, prefix string, depth int) bool {
	paths, err := utils.FindAllCoverageFiles(basePath)
	if err != nil {
		log.Panicf("error file finding coverage files %s", err)
	}

	pkgs, coverageData := getCoverageData(paths, exclusionsMatcher)
	return printCoverage(writer, pkgs, coverageData, float64(minCoverage), prefix, depth)
}

// PrintCoverageSingle is the same as PrintCoverage only for 1 directory only
func PrintCoverageSingle(writer io.Writer, path string, minCoverage int, prefix string, depth int) bool {
	var fullPath string
	if path == "./" {
		fullPath = utils.GetCurrentDir()
	} else {
		fullPath = utils.GetCurrentDir() + path + "/"
	}
	fullPath += "profile.cov"

	contents := getFileContents(fullPath)
	pkgs, coverageData := getCoverageByContents(contents)

	return printCoverage(writer, pkgs, coverageData, float64(minCoverage), prefix, depth)
}

func printCoverage(writer io.Writer, pkgs []string, coverageData coverageByPackage, minCoverage float64, prefix string, depth int) bool {
	addLine(writer)
	fmt.Fprintf(writer, header1Template, "")
	fmt.Fprintf(writer, header2Template, "Cov%", "Cov", "Stmts", "Cov%", "Cov", "Stmts", "Package")
	addLine(writer)

	coverageOk := true
	for _, pkg := range pkgs {
		cover := coverageData[pkg]

		pkgFormatted := strings.Replace(pkg, prefix, "", -1)
		pkgDepth := strings.Count(pkgFormatted, "/")

		if depth > 0 {
			if pkgDepth <= depth {
				if !addLinePrint(writer, pkgFormatted, cover, minCoverage) {
					coverageOk = false
				}
			}
		} else {
			if !addLinePrint(writer, pkgFormatted, cover, minCoverage) {
				coverageOk = false
			}
		}
	}
	addLine(writer)

	return coverageOk
}

func addLine(writer io.Writer) {
	for x := 0; x < 120; x++ {
		fmt.Fprint(writer, "-")
	}
	fmt.Fprint(writer, "\n")
}

const (
	header1Template = "|         Branch        |          Dir          | %-80s |\n"
	header2Template = "| %6s | %6s | %6s | %6s | %6s | %6s | %-80s |\n"
	lineTemplate    = "| %6.2f | %6.0f | %6.0f | %6.2f | %6.0f | %6.0f | %-80s |\n"
	errStart        = "\033[1;31m"
	errEnd          = "\033[0m"
)

func addLinePrint(writer io.Writer, pkgFormatted string, cover *coverage, minCoverage float64) bool {
	precCov, sumCoverage, sumStmts := getSummaryValues(cover)
	precSelf, selfCoverage, selfStmts := getSelfValues(cover)

	if sumCoverage < minCoverage {
		fmt.Fprintf(writer, errStart+lineTemplate+errEnd, precCov, sumCoverage, sumStmts, precSelf, selfCoverage, selfStmts, pkgFormatted)
		return false
	}

	fmt.Fprintf(writer, lineTemplate, precCov, sumCoverage, sumStmts, precSelf, selfCoverage, selfStmts, pkgFormatted)
	return true
}
