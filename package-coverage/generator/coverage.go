package generator

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"io"
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
func generateCoverage(path string, exclusionsMatcher *regexp.Regexp, quiet bool, goTestArgs []string) {
	packageName := findPackageName(path)

	fakeTestFile := addFakeTest(path, packageName)
	defer removeFakeTest(fakeTestFile)

	err := execCoverage(path, coverageFilename, quiet, goTestArgs)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error generating coverage %s", err)
	}

	if exclusionsMatcher == nil {
		return
	}

	err = filterCoverage(filepath.Join(path, coverageFilename), exclusionsMatcher)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error filtering files: %s", err)
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
		utils.LogWhenVerbose("[coverage] err while parsing the '%s' into Go AST Err: '%s", path, err)
		return UnknownPackage
	}

	for pkgName := range pkgs {
		return pkgName
	}
	return UnknownPackage
}

// create a fake test so that all directories are guaranteed to contain tests (and therefore coverage will be generated)
func createTestFile(packageName string, testFilename string) {
	utils.LogWhenVerbose("[coverage] created test for package %s file @ %s", packageName, testFilename)

	if _, err := os.Stat(testFilename); err == nil {
		utils.LogWhenVerbose("[coverage] file already exists @ %s cowardly refusing to overwrite", testFilename)
		return
	}

	file, err := os.OpenFile(testFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while creating test file %s", err)
		return
	}

	_, err = file.WriteString(`package ` + packageName + `

import "testing"

func TestThisTestDoesntReallyTestAnything(t *testing.T) {}
`)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while writing test file %s", err)
		return
	}

	err = file.Close()
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while closing file '%s", err)
	}
}

// remove the previously added fake test (i.e. clean up)
func removeFakeTest(filename string) {
	utils.LogWhenVerbose("[coverage] remove test file @ %s", filename)

	err := os.Remove(filename)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while removing test file @ %s, err: %s", filename, err)
	}
}

// essentially call `go test` to generate the coverage
func execCoverage(dir, coverageFilename string, quiet bool, goTestArgs []string) error {
	command := "go"
	arguments := []string{
		"test",
		"-coverprofile=" + coverageFilename,
	}

	arguments = append(arguments, goTestArgs...)

	cmd := exec.Command(command, arguments...)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while running go test. err: %s", err)
		return err
	}

	utils.LogWhenVerbose("[coverage] created coverage file @ %s%s", dir, coverageFilename)

	if !quiet {
		print(string(out))
	}

	return nil
}

func filterCoverage(coverageFilename string, exclusionsMatcher *regexp.Regexp) error {
	coverageTempFilename := coverageFilename + "~"

	coverageTempFile, err := os.OpenFile(coverageTempFilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while opening file %s, err: %s", coverageTempFilename, err)
		return err
	}

	defer func() {
		if err := coverageTempFile.Close(); err != nil {
			utils.LogWhenVerbose("[coverage] cannot close temporary file %s: %s", coverageTempFilename, err)
		}
	}()

	coverageFile, err := os.OpenFile(coverageFilename, os.O_RDWR, 0)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while openning file %s, err: %s", coverageFilename, err)
		return err
	}

	defer func() {
		if err := coverageFile.Close(); err != nil {
			utils.LogWhenVerbose("[coverage] Cannot close coverage file %s: %s", coverageFilename, err)
		}
	}()

	if err := filterCoverageContents(exclusionsMatcher, coverageFile, coverageTempFile); err != nil {
		utils.LogWhenVerbose("[coverage] error while filtering coverage file %s: %s", coverageFilename, err)
		return err
	}

	if err := os.Remove(coverageFilename); err != nil {
		utils.LogWhenVerbose("[coverage] cannot remove old coverage file %s: %s", coverageFilename, err)
	}

	if err := os.Rename(coverageTempFilename, coverageFilename); err != nil {
		utils.LogWhenVerbose("[coverage] cannot rename filtered coverage file %s: %s", coverageTempFilename, err)
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
