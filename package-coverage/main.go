package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %s\n", r)
		}
	}()

	filename := getFilename()

	contents := getFileContents(filename)

	coverageData := parseLines(contents)

	pkgs := getSortedPackages(coverageData)

	printCoverage(pkgs, coverageData)
}

func getFilename() string {
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		panic("Usage: package-coverage coverage.out")
	}

	return filename
}

func getFileContents(filename string) string {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(contents)
}

func getSortedPackages(coverageData map[string]*coverage) []string {
	output := []string{}

	for pkg := range coverageData {
		output = append(output, pkg)
	}

	sort.Strings(output)

	return output
}

func printCoverage(pkgs []string, coverageData map[string]*coverage) {
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
