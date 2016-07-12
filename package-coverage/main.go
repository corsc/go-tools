package main

import (
	"flag"
	"fmt"
	"regexp"

	"os"

	"github.com/corsc/go-tools/package-coverage/generator"
	"github.com/corsc/go-tools/package-coverage/parser"
	"github.com/corsc/go-tools/package-coverage/utils"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %s\n", r)
		}
	}()

	verbose := false
	coverage := false
	singleDir := false
	clean := false
	print := false
	ignoreDirs := ""
	minCoverage := 0
	var matcher *regexp.Regexp

	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&coverage, "c", false, "generate coverage")
	flag.BoolVar(&singleDir, "s", false, "only generate for the supplied directory (no recursion)")
	flag.BoolVar(&clean, "d", false, "clean")
	flag.BoolVar(&print, "p", false, "print coverage to stdout")
	flag.StringVar(&ignoreDirs, "i", `./\.git.*|./_.*`, "ignore regex specified directory")
	flag.IntVar(&minCoverage, "m", 0, "minimum coverage")
	flag.Parse()

	if !verbose {
		utils.VerboseOff()
	}

	startDir := utils.GetCurrentDir()
	path := getPath()

	if ignoreDirs != "" {
		matcher = regexp.MustCompile(ignoreDirs)
	}

	if coverage {
		if singleDir {
			generator.CoverageSingle(path, matcher)
		} else {
			generator.Coverage(path, matcher)
		}
	}

	// switch back to start dir
	err := os.Chdir(startDir)
	if err != nil {
		panic(err)
	}

	var coverageOk bool
	if print {
		if singleDir {
			coverageOk = parser.PrintCoverageSingle(path, matcher, minCoverage)
		} else {
			coverageOk = parser.PrintCoverage(path, matcher, minCoverage)
		}
	}

	if clean {
		generator.Clean(path, matcher)
	}

	if !coverageOk {
		os.Exit(-1)
	}
}

func getPath() string {
	path := flag.Arg(0)
	if path == "" {
		panic("Please include a directory as the last argument")
	}
	return path
}
