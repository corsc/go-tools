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

	"github.com/corsc/go-tools/commons"
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
func generateCoverage(path string, exclusions *regexp.Regexp, quietMode bool, testArgs []string) {
	packageName := findPackageName(path)

	fakeTestFile := addFakeTest(path, packageName)
	defer removeFakeTest(fakeTestFile)

	err := execCoverage(path, quietMode, testArgs)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error generating coverage %s", err)
	}

	if exclusions == nil {
		return
	}

	err = filterCoverage(filepath.Join(path, coverageFilename), exclusions)
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

	file, err := os.OpenFile(testFilename, os.O_RDWR|os.O_CREATE, 0600)
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
func execCoverage(dir string, quiet bool, goTestArgs []string) error {
	arguments := []string{
		"test",
		"-coverprofile=" + coverageFilename,
	}

	arguments = append(arguments, goTestArgs...)

	cmd := exec.Command("go", arguments...)
	cmd.Dir = dir

	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while running go test. err: %s", err)
		return err
	}

	utils.LogWhenVerbose("[coverage] created coverage file @ %s%s", dir, coverageFilename)

	return nil
}

func filterCoverage(coverageFilename string, exclusionsMatcher *regexp.Regexp) error {
	coverageTempFilename := coverageFilename + "~"

	coverageTempFile, err := os.OpenFile(coverageTempFilename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while opening file %s, err: %s", coverageTempFilename, err)
		return err
	}

	defer commons.CloseIO(coverageTempFile)

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
			utils.LogWhenVerbose("[coverage] skipped file %s", fileName)
			continue
		}

		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}

	return nil
}
