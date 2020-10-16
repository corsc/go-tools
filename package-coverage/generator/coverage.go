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

	"github.com/corsc/go-commons/iocloser"
	"github.com/corsc/go-tools/package-coverage/utils"
)

const coverageFilename = "profile.cov"

var fakeTestFilename = "fake_test.go"
var fakeCodeFilename = "fake_code.go"

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
func generateCoverage(path string, exclusions *regexp.Regexp, quietMode, race bool, tags string) {
	packageName := findPackageName(path)

	addFakes(path, packageName)

	err := execCoverage(path, quietMode, race, tags)
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

	re
}

func addFakes(path, packageName string) {
	testFilename := createTestFilename(path)
	createTestFile(packageName, testFilename)

	codeFilename := createCodeFilename(path)
	createCodeFile(packageName, codeFilename)
}

func removeFakes(path string) {
	testFilename := createTestFilename(path)
	utils.LogWhenVerbose("[coverage] remove fake test @ %s", testFilename)
	err := os.Remove(testFilename)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while removing fake test @ %s, err: %s", testFilename, err)
	}

	codeFilename := createCodeFilename(path)
	utils.LogWhenVerbose("[coverage] remove fake code @ %s", codeFilename)
	err = os.Remove(codeFilename)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while removing fake code @ %s, err: %s", codeFilename, err)
	}
}

func createTestFilename(path string) string {
	return path + fakeTestFilename
}

func createCodeFilename(path string) string {
	return path + fakeCodeFilename
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

// create a fake function so that all directories are guaranteed to contain code (and therefore coverage will be generated)
func createCodeFile(packageName string, codeFilename string) {
	utils.LogWhenVerbose("[coverage] created code for package %s file @ %s", packageName, codeFilename)

	if _, err := os.Stat(codeFilename); err == nil {
		utils.LogWhenVerbose("[coverage] file already exists @ %s cowardly refusing to overwrite", codeFilename)
		return
	}

	file, err := os.OpenFile(codeFilename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while creating code file %s", err)
		return
	}

	_, err = file.WriteString(`package ` + packageName + `

func ThisCodeDoesntReallyDoAnything() {}
`)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while writing code file %s", err)
		return
	}

	err = file.Close()
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while closing file '%s", err)
	}
}

// essentially call `go test` to generate the coverage
func execCoverage(dir string, quiet, race bool, tags string) error {
	arguments := []string{
		"test",
		"-coverprofile=" + coverageFilename,
	}

	if race {
		arguments = append(arguments, `--race`)
	}

	if len(tags) > 0 {
		arguments = append(arguments, `-tags="`+tags+`"`)
	}

	cmd := exec.Command("go", arguments...)
	cmd.Dir = dir

	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	payload, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogAlways("failed to get test command output. dir: %s err: %s", dir, err)
	}

	if err != nil {
		utils.LogAlways("[coverage] test output %s:\n%s", dir, payload)

		if err.Error() == "exit status 1" {
			utils.LogAlways("WARNING: tests for %s are broken. err: %s", dir, err)
		} else {
			utils.LogAlways("[coverage] error while running go test %s. err: %s", dir, err)
		}

		return err
	} else {
		utils.LogWhenVerbose("[coverage] test output %s:\n%s", dir, payload)
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

	defer iocloser.Close(coverageTempFile)

	coverageFile, err := os.OpenFile(coverageFilename, os.O_RDWR, 0)
	if err != nil {
		utils.LogWhenVerbose("[coverage] error while opening file %s, err: %s", coverageFilename, err)
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
		if fileNameEndIndex == -1 {
			utils.LogWhenVerbose("[coverage] error in line '%s'", line)
			continue
		}

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
