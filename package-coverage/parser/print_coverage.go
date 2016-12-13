package parser

import (
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// CoverageByPackage contains the result of parsing one or more package's coverage file
type coverageByPackage map[string]*coverage

// PrintCoverage will attempt to calculate and print the coverage from the supplied coverage file to standard out.
func PrintCoverage(writer io.Writer, basePath string, exclusionsMatcher *regexp.Regexp, minCoverage int) bool {
	paths, err := utils.FindAllCoverageFiles(basePath)
	if err != nil {
		log.Panicf("error file finding coverage files %s", err)
	}

	pkgs, coverageData := getCoverageData(paths, exclusionsMatcher)
	return printCoverage(writer, pkgs, coverageData, float64(minCoverage))
}

// PrintCoverageSingle is the same as PrintCoverage only for 1 directory only
func PrintCoverageSingle(writer io.Writer, path string, minCoverage int) bool {
	var fullPath string
	if path == "./" {
		fullPath = utils.GetCurrentDir()
	} else {
		fullPath = utils.GetCurrentDir() + path + "/"
	}
	fullPath += "profile.cov"

	contents := getFileContents(fullPath)
	pkgs, coverageData := getCoverageByContents(contents)

	return printCoverage(writer, pkgs, coverageData, float64(minCoverage))
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
