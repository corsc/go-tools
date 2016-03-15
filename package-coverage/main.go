package main

import (
	"flag"
	"fmt"
	"regexp"

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

	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&coverage, "c", false, "generate coverage")
	flag.BoolVar(&singleDir, "s", false, "only generate for the supplied directory (no recursion)")
	flag.BoolVar(&clean, "d", false, "clean")
	flag.BoolVar(&print, "p", false, "print coverage to stdout")
	flag.StringVar(&ignoreDirs, "i", `./\.git.*|./_.*`, "ignore regex specified directory")
	flag.Parse()

	if !verbose {
		utils.VerboseOff()
	}

	path := getPath()

	if ignoreDirs != "" {
		if matched, _ := regexp.MatchString(ignoreDirs, path); matched {
			return
		}
	}

	if coverage {
		if singleDir {
			generator.CoverageSingle(path)
		} else {
			generator.Coverage(path)
		}
	}

	if print {
		parser.PrintCoverage(path)
	}

	if clean {
		generator.Clean(path)
	}
}

func getPath() string {
	path := flag.Arg(0)
	if path == "" {
		panic("Please include a directory as the last argument")
	}
	return path
}
