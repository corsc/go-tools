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
	"io/ioutil"
	"regexp"
	"sort"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// get coverage using the paths and exclusions supplied
func getCoverageData(paths []string, exclusionsMatcher *regexp.Regexp) ([]string, coverageByPackage) {
	var contents string
	for _, path := range paths {
		if exclusionsMatcher.FindString(path) != "" {
			utils.LogWhenVerbose("[print] Printing of coverage for path '%s' skipped due to exclusions regex '%s'",
				path, exclusionsMatcher.String())
			continue
		}

		contents += getFileContents(path)
	}

	return getCoverageByContents(contents)
}

// get coverage from supplied string (used after concatenating all the individual coverage files together
func getCoverageByContents(contents string) ([]string, coverageByPackage) {
	coverageData := getCoverageByPackage(contents)
	pkgs := getSortedPackages(coverageData)

	return pkgs, coverageData
}

// will calculate and return the coverage for a package or packages from the supplied coverage file contents
func getCoverageByPackage(contents string) coverageByPackage {
	coverageData := parseLines(contents)
	return coverageData
}

func getFileContents(filename string) string {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(contents)
}

func getSortedPackages(coverageData coverageByPackage) []string {
	output := []string{}

	for pkg := range coverageData {
		output = append(output, pkg)
	}

	sort.Strings(output)

	return output
}

func getSummaryValues(cover *coverage) (float64, float64, float64) {
	stmts := float64(cover.selfStatements + cover.childStatements)
	stmtsCovered := float64(cover.selfCovered + cover.childCovered)

	perc := getPercentage(stmts, stmtsCovered)
	return perc, stmtsCovered, stmts
}

func getSelfValues(cover *coverage) (float64, float64, float64) {
	stmts := float64(cover.selfStatements)
	stmtsCovered := float64(cover.selfCovered)

	perc := getPercentage(stmts, stmtsCovered)
	return perc, stmtsCovered, stmts
}

func getPercentage(stmts float64, stmtsCovered float64) float64 {
	if stmts <= 0.0 {
		return 100
	}

	covered := (stmtsCovered / stmts) * 100

	return covered
}
