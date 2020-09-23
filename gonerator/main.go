// Copyright 2017 Corey Scott http://www.sage42.org/
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

// Gonerator is a tool to automate the creation of go code from inputs and templates
// Given the name of a type T, this tool will create a new self-contained Go source file.
//
// The file is created in the same package and directory as the package that defines T.
// If the file exists the generation will not occur.
//
// To use this tool install it with go get github.com/corsc/go-tools/gonerator
//
// Run this using the command:
//
//	gonerator -type=MyType -i=template.tmpl -o my_type.go
//
// will create the files my_type.go
//
// This code was adapted and extended from https://github.com/golang/tools/blob/master/cmd/stringer/stringer.go

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/corsc/go-tools/gonerator/gonerator"
)

func main() {
	setUp()
	getInputs()

	dir := "./"
	g := &gonerator.Gonerator{}
	g.ParsePackageDir(dir)
	g.Build(dir, *typeName, *templateFile, *outputFile, *extras, *dryRun, *noop)
}

var (
	typeName     = flag.String("i", "", "type name; must be set")
	templateFile = flag.String("t", "", "template file; one of template and noop must be set")
	noop         = flag.Bool("noop", false, "generate a NO-OP implementation of the supplied interface; one of template and noop must be set")
	outputFile   = flag.String("o", "", "output file; must be set")
	dryRun       = flag.Bool("d", false, "dry-run; output to stdOut instead of updating the file")
	extras       = flag.String("e", "", "comma separated list of extra values; optional")
)

// Usage outputs the usage of this tool to std err
func Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "\tgonerator [flags] -i=MyType -t=template.tmpl [-o=mytype_gonerated.go]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\t\tor\n")
	_, _ = fmt.Fprintf(os.Stderr, "\tgonerator [flags] -i=MyType -noop [-o=mytype_gonerated.go]\n")
	_, _ = fmt.Fprintf(os.Stderr, "\t\tor\n")
	_, _ = fmt.Fprintf(os.Stderr, "\t//go:generate gonerator [flags] -i=MyType -t=template.tmpl -o=mytype_gonerated.go]\n")
	_, _ = fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func setUp() {
	log.SetFlags(0)
	log.SetPrefix("gonerator: ")
}

func getInputs() {
	flag.Usage = Usage
	flag.Parse()

	if len(*typeName) == 0 || len(*outputFile) == 0 || (len(*templateFile) == 0 && *noop == false) {
		flag.Usage()
		os.Exit(2)
	}
}
