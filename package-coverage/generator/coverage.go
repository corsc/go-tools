package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/corsc/go-tools/package-coverage/utils"
)

const coverageFilename = "profile.cov"

var fakeTestFilename = "fake_test.go"

func processAllDirs(basePath string, exclusionsMatcher *regexp.Regexp, logTag string, actionFunc func(string)) {
	paths, err := utils.FindAllGoDirs(basePath)
	if err != nil {
		return
	}

	for _, path := range paths {
		if exclusionsMatcher.FindString(path) != "" {
			utils.LogWhenVerbose("[%s] path '%s' skipped due to skipDir regex '%s'",
				logTag, path, exclusionsMatcher.String())
			continue
		}

		utils.LogWhenVerbose("[%s] processing path '%s'", logTag, path)
		actionFunc(path)
	}
}

// this function will generate the test coverage for the supplied directory
func generateCoverage(path string, exclusionsMatcher *regexp.Regexp, goTestArgs []string) {
	packageName := findPackageName(path)

	fakeTestFile := addFakeTest(path, packageName)
	defer removeFakeTest(fakeTestFile)

	err := execCoverage(path, coverageFilename, goTestArgs)
	if err != nil {
		log.Printf("error generating coverage %s", err)
	}

	if exclusionsMatcher == nil {
		return
	}

	err = filterCoverage(filepath.Join(path, coverageFilename), exclusionsMatcher)
	if err != nil {
		log.Printf("error filtering files: %s", err)
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

func filterCoverage(coverageFilename string, exclusionsMatcher *regexp.Regexp) error {
	coverageTempFilename := coverageFilename + "~"

	coverageTempFile, err := os.OpenFile(coverageTempFilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer func() {
		if err := coverageTempFile.Close(); err != nil {
			panic(fmt.Errorf("cannot close temporary file %s: %v", coverageTempFilename, err))
		}
	}()

	coverageFile, err := os.OpenFile(coverageFilename, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	defer func() {
		if err := coverageFile.Close(); err != nil {
			panic(fmt.Errorf("Cannot close coverage file %s: %v", coverageFilename, err))
		}
	}()

	if err := filterCoverageContents(exclusionsMatcher, coverageFile, coverageTempFile); err != nil {
		return err
	}

	if err := os.Remove(coverageFilename); err != nil {
		log.Printf("cannot remove old coverage file %s: %v", coverageFilename, err)
	}

	if err := os.Rename(coverageTempFilename, coverageFilename); err != nil {
		log.Printf("cannot rename filtered coverage file %s: %v", coverageTempFilename, err)
	}

	return nil
}

func filterCoverageContents(exclusionsMatcher *regexp.Regexp, in io.Reader, out io.Writer) error {
	coverageFileScanner := bufio.NewScanner(in)
	for coverageFileScanner.Scan() {
		line := coverageFileScanner.Text()

		fileNameEndIndex := strings.LastIndex(line, ":")
		fileName := line[:fileNameEndIndex]

		if exclusionsMatcher.MatchString(fileName) {
			continue
		}

		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}

	return nil
}
