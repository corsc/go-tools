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
//	gonerator -type=MyType -template=template.tmpl -test-template=test-template.tmpl
//
// will create the files mytype_gonerated.go & mytype_gonerated_test.go
//
// Typically this process would be run using go generate, like this:
//
//	// go:generate gonerator -type=MyType -template=template.tmpl -test-template=test-template.tmpl
//
// This code was adapted and extended from https://github.com/golang/tools/blob/master/cmd/stringer/stringer.go
package main

import (
	"flag"
	"log"
	"os"

	"github.com/corsc/go-tools/commons"
	"github.com/corsc/go-tools/gonerator/gonerator"
)

func main() {
	setUp()
	typeName, templateFile, outputFile, extras := getInputs()

	dir := "./"
	g := &gonerator.Gonerator{}
	g.ParsePackageDir(dir)
	g.Build(dir, typeName, templateFile, outputFile, extras)
}

var (
	typeName     = flag.String("i", "", "type name; must be set")
	templateFile = flag.String("t", "", "template file; must be set")
	outputFile   = flag.String("o", "", "output file; must be set")
	extras       = flag.String("e", "", "comma separated list of extra values; optional")
)

// Usage outputs the usage of this tool to std err
func Usage() {
	commons.LogError("Usage of %s:\n", os.Args[0])
	commons.LogError("\tgonerator [flags] -i=T -t=template.tmpl [-o-T_gonerated.go]\n")
	commons.LogError("Flags:\n")
	flag.PrintDefaults()
}

func setUp() {
	log.SetFlags(0)
	log.SetPrefix("gonerator: ")
}

func getInputs() (string, string, string, string) {
	flag.Usage = Usage
	flag.Parse()

	if len(*typeName) == 0 || len(*templateFile) == 0 || len(*outputFile) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	return *typeName, *templateFile, *outputFile, *extras
}
