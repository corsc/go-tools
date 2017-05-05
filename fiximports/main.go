// Package main is the main package for fix imports
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/corsc/go-tools/commons"
	"github.com/corsc/go-tools/fiximports/fiximports"
)

func usage() {
	commons.LogError("Usage of %s:\n", os.Args[0])
	commons.LogError("\tfiximports [flags] # runs on package in current directory\n")
	commons.LogError("\tfiximports [flags] directory\n")
	commons.LogError("\tfiximports [flags] files... # must be a single package\n")
	commons.LogError("Flags:\n")
	flag.PrintDefaults()
}

func main() {
	updateFile := false

	flag.Usage = usage
	flag.BoolVar(&updateFile, "w", false, "write result to (source) file instead of stdout")
	flag.Parse()

	var filenames []string
	var err error

	switch flag.NArg() {
	case 0:
		filenames, err = commons.GetGoFilesFromCurrentDir()

	case 1:
		arg := flag.Arg(0)
		if strings.HasSuffix(arg, "/...") && commons.IsDir(arg[:len(arg)-4]) {
			filenames, err = commons.GetGoFilesFromDirectoryRecursive(arg)

		} else if commons.IsDir(arg) {
			filenames, err = commons.GetGoFilesFromDir(arg)

		} else if commons.FileExists(arg) {
			filenames, err = commons.GetGoFiles(arg)

		} else {
			err = fmt.Errorf("'%s' did not resolve to a directory or file", arg)
		}

	default:
		filenames, err = commons.GetGoFiles(flag.Args()...)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	var outputWriter io.Writer
	if !updateFile {
		// default write to os.Stdout
		outputWriter = os.Stdout
	}

	fiximports.ProcessFiles(filenames, outputWriter)
}
