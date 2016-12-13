package parser

import (
	"io/ioutil"
	"log"
	"regexp"
	"sort"
)

// get coverage using the paths and exclusions supplied
func getCoverageData(paths []string, exclusionsMatcher *regexp.Regexp) ([]string, coverageByPackage) {
	var contents string
	for _, path := range paths {
		if exclusionsMatcher.FindString(path) != "" {
			log.Printf("Printing of coverage for path '%s' skipped due to exclusions regex '%s'",
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

// calculate the coverage and statement counts from the supplied data
func getStats(cover *coverage) (float64, float64) {
	stmts := float64(cover.selfStatements + cover.childStatements)
	stmtsCovered := float64(cover.selfCovered + cover.childCovered)

	covered := (stmtsCovered / stmts) * 100

	return covered, stmts
}
