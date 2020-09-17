// Copyright 2017- Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main is the main package for fix imports
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

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

	argsToFile := fiximports.FilesFromArgsFactory(flag.NArg())
	filenames, err := argsToFile.FileNames()

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
