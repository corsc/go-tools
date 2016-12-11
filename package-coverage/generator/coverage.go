package generator

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"

	"bytes"

	"regexp"

	"github.com/corsc/go-tools/package-coverage/utils"
)

const coverageFilename = "profile.cov"

var fakeTestFilename = "fake_test.go"

func processAllDirs(basePath string, matcher *regexp.Regexp, logTag string, actionFunc func(string)) {
	paths, err := utils.FindAllGoDirs(basePath)
	if err != nil {
		return
	}

	for _, path := range paths {
		if matcher.FindString(path) != "" {
			utils.LogWhenVerbose("[%s] path '%s' skipped due to skipDir regex '%s'", logTag, path, matcher.String())
			continue
		}

		utils.LogWhenVerbose("[%s] processing path '%s'", logTag, path)
		actionFunc(path)
	}
}

// this function will cause the generation of test coverage for the supplied directory and return the file path of the
// resultant coverage file
func generateCoverage(path string, goTestArgs []string) {
	packageName := findPackageName(path)

	fakeTestFile := addFakeTest(path, packageName)
	defer removeFakeTest(fakeTestFile)

	err := execCoverage(path, coverageFilename, goTestArgs)
	if err != nil {
		log.Printf("error generating coverage %s", err)
	}
}

// add a fake test to ensure that there is at least 1 test in this directory
func addFakeTest(path string, packageName string) string {
	testFilename := createTestFilename(path)

	createTestFile(packageName, testFilename)

	return testFilename
}

func createTestFilename(path string) string {
	return path + fakeTestFilename
}

// find the package name by using the go AST
func findPackageName(path string) string {
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, path, nil, 0)
	if err != nil {
		log.Printf("err while parsing the '%s' into Go AST Err: '%s", path, err)
		return UnknownPackage
	}

	for pkgName := range pkgs {
		return pkgName
	}
	return UnknownPackage
}

// create a fake test so that all directories are guaranteed to contain tests (and therefore coverage will be generated)
func createTestFile(packageName string, testFilename string) {
	utils.LogWhenVerbose("created test for package %s file @ %s", packageName, testFilename)

	if _, err := os.Stat(testFilename); err == nil {
		utils.LogWhenVerbose("file already exists @ %s cowardly refusing to overwrite", testFilename)
		return
	}

	file, err := os.OpenFile(testFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("error while creating test file %s", err)
		return
	}

	_, err = file.WriteString(`package ` + packageName + `

import "testing"

func TestThisTestDoesntReallyTestAnything(t *testing.T) {}
`)
	if err != nil {
		log.Printf("error while writing test file %s", err)
		return
	}

	err = file.Close()
	if err != nil {
		log.Printf("error while closing file '%s", err)
	}
}

// remove the previously added fake test (i.e. clean up)
func removeFakeTest(filename string) {
	utils.LogWhenVerbose("remove test file @ %s", filename)

	err := os.Remove(filename)
	if err != nil {
		log.Printf("error while removing test file @ %s, err: %s", filename, err)
	}
}

// essentially call `go test` to generate the coverage
var execCoverage = func(dir, coverageFilename string, goTestArgs []string) error {
	var stdErr bytes.Buffer

	command := "go"
	arguments := []string{
		"test",
		"-coverprofile=" + coverageFilename,
	}

	arguments = append(arguments, goTestArgs...)

	cmd := exec.Command(command, arguments...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		log.Printf("error while running go test. stdErr: %s\nerr: %s", stdErr.String(), err)
		return err
	}
	utils.LogWhenVerbose("created coverage file @ %s%s", dir, coverageFilename)
	return nil
}
