package main

import (
	"flag"
	"fmt"

	"sage42.org/go-tools/package-coverage/generator"
	"sage42.org/go-tools/package-coverage/utils"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %s\n", r)
		}
	}()

	verbose := false
	coverage := false
	clean := false
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.BoolVar(&coverage, "c", false, "generate coverage")
	flag.BoolVar(&clean, "d", false, "clean")
	flag.Parse()

	if !verbose {
		utils.VerboseOff()
	}

	path := getPath()
	if coverage {
		generator.Coverage(path)
	}

	if clean {
		generator.Clean(path)
	}
}

func getPath() string {
	path := flag.Arg(0)
	if path == "" {
		panic("Usage: package-coverage ./")
	}
	return path
}