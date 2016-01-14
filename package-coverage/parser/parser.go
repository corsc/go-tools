package parser

import (
	"fmt"
	"io/ioutil"
	"sort"
)

// CoverageByPackage contains the result of parsing one or more package's coverage file
type CoverageByPackage map[string]*coverage

// PrintCoverage will attempt to calculate and print the coverage from the supplied coverage file
func PrintCoverage(filename string) {
	coverageData := CalculateCoverage(filename)

	pkgs := getSortedPackages(coverageData)

	printCoverage(pkgs, coverageData)
}

// CalculateCoverage will calculate and return the coverage for a package or packages from the supplied file
func CalculateCoverage(filename string) CoverageByPackage {
	contents := getFileContents(filename)

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

func getSortedPackages(coverageData CoverageByPackage) []string {
	output := []string{}

	for pkg := range coverageData {
		output = append(output, pkg)
	}

	sort.Strings(output)

	return output
}

func printCoverage(pkgs []string, coverageData CoverageByPackage) {
	fmt.Printf("  %%		Statements	Package\n")

	for _, pkg := range pkgs {
		cover := coverageData[pkg]

		stmts := float64(cover.selfStatements + cover.childStatements)
		stmtsCovered := float64(cover.selfCovered + cover.childCovered)

		covered := (stmtsCovered / stmts) * 100

		fmt.Printf("%2.2f		%5.0f		%s\n", covered, stmts, pkg)
	}
	fmt.Println()
}
