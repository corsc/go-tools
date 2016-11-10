package parser

import (
	"fmt"
	"io/ioutil"
	"sort"

	"log"
	"regexp"

	"io"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// CoverageByPackage contains the result of parsing one or more package's coverage file
type coverageByPackage map[string]*coverage

// PrintCoverageSingle is the same as PrintCoverage only for 1 directory only
func PrintCoverageSingle(writer io.Writer, path string, matcher *regexp.Regexp) {
	var fullPath string
	if path == "./" {
		fullPath = utils.GetCurrentDir()
	} else {
		fullPath = utils.GetCurrentDir() + path + "/"
	}
	fullPath += "profile.cov"

	print(writer, []string{fullPath}, matcher)
}

// PrintCoverage will attempt to calculate and print the coverage from the supplied coverage file
func PrintCoverage(writer io.Writer, basePath string, matcher *regexp.Regexp) {
	paths, err := utils.FindAllCoverageFiles(basePath)
	if err != nil {
		log.Panicf("error file finding coverage files %s", err)
	}
	print(writer, paths, matcher)
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

func print(writer io.Writer, paths []string, matcher *regexp.Regexp) {
	pkgs, coverageData := getCoverageData(paths, matcher)
	printCoverage(writer, pkgs, coverageData)
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

func printCoverage(writer io.Writer, pkgs []string, coverageData coverageByPackage) {
	fmt.Fprintf(writer, "  %%		Statements	Package\n")

	for _, pkg := range pkgs {
		cover := coverageData[pkg]
		covered, stmts := getStats(cover)

		fmt.Fprintf(writer, "%2.2f		%5.0f		%s\n", covered, stmts, pkg)
	}
	fmt.Fprint(writer, "\n")
}

func getStats(cover *coverage) (float64, float64) {
	stmts := float64(cover.selfStatements + cover.childStatements)
	stmtsCovered := float64(cover.selfCovered + cover.childCovered)

	covered := (stmtsCovered / stmts) * 100

	return covered, stmts
}
