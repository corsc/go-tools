package parser

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"sort"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// CoverageByPackage contains the result of parsing one or more package's coverage file
type coverageByPackage map[string]*coverage

// PrintCoverageSingle is the same as PrintCoverage only for 1 directory only
func PrintCoverageSingle(writer io.Writer, path string, matcher *regexp.Regexp, minCoverage int) bool {
	var fullPath string
	if path == "./" {
		fullPath = utils.GetCurrentDir()
	} else {
		fullPath = utils.GetCurrentDir() + path + "/"
	}
	fullPath += "profile.cov"

	return print(writer, []string{fullPath}, matcher, minCoverage)
}

// PrintCoverage will attempt to calculate and print the coverage from the supplied coverage file
func PrintCoverage(writer io.Writer, basePath string, matcher *regexp.Regexp, minCoverage int) bool {
	paths, err := utils.FindAllCoverageFiles(basePath)
	if err != nil {
		log.Panicf("error file finding coverage files %s", err)
	}
	return print(writer, paths, matcher, minCoverage)
}

func getCoverageData(paths []string, matcher *regexp.Regexp) ([]string, coverageByPackage) {
	var contents string
	for _, path := range paths {
		if matcher.FindString(path) != "" {
			log.Printf("Printing of coverage for path '%s' skipped due to skipDir regex '%s'", path, matcher.String())
			continue
		}

		contents += getFileContents(path)
	}

	coverageData := calculateCoverage(contents)
	pkgs := getSortedPackages(coverageData)

	return pkgs, coverageData
}

func print(writer io.Writer, paths []string, matcher *regexp.Regexp, minCoverage int) bool {
	pkgs, coverageData := getCoverageData(paths, matcher)
	return printCoverage(writer, pkgs, coverageData, float64(minCoverage))
}

// CalculateCoverage will calculate and return the coverage for a package or packages from the supplied coverage file contents
func calculateCoverage(contents string) coverageByPackage {
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

func printCoverage(writer io.Writer, pkgs []string, coverageData coverageByPackage, minCoverage float64) bool {
	fmt.Fprintf(writer, "  %%		Statements	Package\n")

	coverageOk := true
	for _, pkg := range pkgs {
		cover := coverageData[pkg]
		covered, stmts := getStats(cover)

		if covered < minCoverage {
			fmt.Fprintf(writer, "\033[1;31m%2.2f		%5.0f		%s\033[0m\n", covered, stmts, pkg)
			coverageOk = false
		} else {
			fmt.Fprintf(writer, "%2.2f		%5.0f		%s\n", covered, stmts, pkg)
		}
	}
	fmt.Fprint(writer, "\n")

	return coverageOk
}

func getStats(cover *coverage) (float64, float64) {
	stmts := float64(cover.selfStatements + cover.childStatements)
	stmtsCovered := float64(cover.selfCovered + cover.childCovered)

	covered := (stmtsCovered / stmts) * 100

	return covered, stmts
}
