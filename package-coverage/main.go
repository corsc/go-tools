package main

import (
	"flag"
	"fmt"
	"regexp"

	"os"

	"bytes"

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
	slack := false
	ignoreDirs := ""
	webHook := ""
	prefix := ""
	depth := 0
	minCoverage := 0
	var matcher *regexp.Regexp

	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&coverage, "c", false, "generate coverage")
	flag.BoolVar(&singleDir, "s", false, "only generate for the supplied directory (no recursion / will ignore -i)")
	flag.BoolVar(&clean, "d", false, "clean")
	flag.BoolVar(&print, "p", false, "print coverage to stdout")
	flag.BoolVar(&slack, "slack", false, "output coverage to slack")
	flag.StringVar(&ignoreDirs, "i", `./\.git.*|./_.*`, "ignore regex specified directory")
	flag.StringVar(&webHook, "webhook", "", "Slack webhook URL")
	flag.StringVar(&prefix, "prefix", "", "Prefix to be removed from the output (currently only supported by Slack output)")
	flag.IntVar(&depth, "depth", 0, "How many levels of coverage to output (default is 0 = all) (currently only supported by Slack output)")
	flag.IntVar(&minCoverage, "m", 0, "minimum coverage")
	flag.Parse()

	if !verbose {
		utils.VerboseOff()
	}

	startDir := utils.GetCurrentDir()
	path := getPath()
	goTestArgs := getGoTestArguments()

	if ignoreDirs != "" {
		matcher = regexp.MustCompile(ignoreDirs)
	}

	if coverage {
		if singleDir {
			generator.CoverageSingle(path, verbose, goTestArgs)
		} else {
			generator.Coverage(path, matcher, verbose, goTestArgs)
		}
	}

	// switch back to start dir
	err := os.Chdir(startDir)
	if err != nil {
		panic(err)
	}

	coverageOk := true
	if print {
		buffer := bytes.Buffer{}

		if singleDir {
			coverageOk = parser.PrintCoverageSingle(&buffer, path, minCoverage)
		} else {
			coverageOk = parser.PrintCoverage(&buffer, path, matcher, minCoverage)
		}

		fmt.Print(buffer.String())
	}

	if slack {
		if singleDir {
			parser.SlackCoverageSingle(path, webHook, prefix, depth)
		} else {
			parser.SlackCoverage(path, matcher, webHook, prefix, depth)
		}
	}

	if clean {
		if singleDir {
			generator.CleanSingle(path)
		} else {
			generator.Clean(path, matcher)
		}
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

func getGoTestArguments() []string {
	args := flag.Args()

	// We only assume what comes after -- to be `go test` arguments. If there are two arguments, we do not assume them
	// to be `go test` arguments.
	if (len(args) >= 2 && args[1] != "--") || len(args) < 3 {
		return []string{}
	}

	return args[2:]
}
