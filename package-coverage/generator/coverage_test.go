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
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/corsc/go-tools/package-coverage/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	// change test file so that it doesn't overlap when using this tool on itself
	fakeTestFilename = "my-fake-test.go"
}

func TestProcessAllDirs(t *testing.T) {
	dir := strings.TrimSuffix(utils.GetCurrentDir(), "/generator/")
	tests := []struct {
		Name    string
		Exclude string
		Result  []string
	}{
		{
			Name:    "excludesDir",
			Exclude: `/excluded/`,
			Result: []string{
				dir + "/test-data/pathmatcher/",
				dir + "/test-data/pathmatcher/included/",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			matcher := regexp.MustCompile(test.Exclude)
			var matched []string
			processAllDirs("../test-data/pathmatcher", matcher, "processAllDirs", func(path string) {
				matched = append(matched, path)
			})

			sort.Strings(test.Result)
			sort.Strings(matched)
			assert.Equal(t, test.Result, matched, "expected processed files to match")
		})
	}

}

func TestAddFakes_HappyPath(t *testing.T) {
	path := utils.GetCurrentDir()
	packageName := "generator"

	expectedTestFilename := path + fakeTestFilename
	defer removeTestFile(expectedTestFilename)

	expectedCodeFilename := path + fakeCodeFilename
	defer removeTestFile(expectedCodeFilename)

	addFakes(path, packageName)

	assertFileExists(t, expectedTestFilename)
	assertFileExists(t, expectedCodeFilename)
}

func TestCreateTestFilename(t *testing.T) {
	path := utils.GetCurrentDir() + "generator/"
	expected := path + fakeTestFilename

	result := createTestFilename(path)
	assert.Equal(t, expected, result)
}

func TestExtractPackageName(t *testing.T) {
	path := utils.GetCurrentDir()
	expected := "generator"

	result := findPackageName(path)
	assert.Equal(t, expected, result)
}

func TestCreateTestFile(t *testing.T) {
	packageName := "mypackage"
	testFile := utils.GetCurrentDir() + "my_test.go"

	defer removeTestFile(testFile)

	assertFileNotExists(t, testFile)

	createTestFile(packageName, testFile)

	assertFileExists(t, testFile)
}

func removeTestFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Panicf("error while removing test file %s", err)
	}
}

func assertFileExists(t *testing.T, path string) {
	_, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	assert.True(t, os.IsExist(err))
}

func assertFileNotExists(t *testing.T, path string) {
	_, err := os.OpenFile(path, os.O_RDONLY, 0666)
	assert.True(t, os.IsNotExist(err))
}

func TestFilterCoverage(t *testing.T) {
	const (
		sampleCoverageFile = `github.com/corsc/go-tools/package-coverage/generator/generator.go:9.87,10.76 1 0
github.com/corsc/go-tools/package-coverage/generator/generator.go:10.76,12.3 1 0
github.com/corsc/go-tools/package-coverage/generator/generator.go:16.93,18.2 1 0
github.com/corsc/go-tools/package-coverage/generator/clean.go:13.63,15.2 1 0
github.com/corsc/go-tools/package-coverage/generator/clean.go:18.31,20.2 1 0
github.com/corsc/go-tools/package-coverage/generator/clean.go:22.25,25.2 2 0
github.com/corsc/go-tools/package-coverage/generator/clean.go:28.42,29.45 1 0
github.com/corsc/go-tools/package-coverage/generator/clean.go:29.45,33.17 3 0
github.com/corsc/go-tools/package-coverage/generator/clean.go:33.17,35.4 1 0
`
		filteredCoverageFile = `github.com/corsc/go-tools/package-coverage/generator/generator.go:9.87,10.76 1 0
github.com/corsc/go-tools/package-coverage/generator/generator.go:10.76,12.3 1 0
github.com/corsc/go-tools/package-coverage/generator/generator.go:16.93,18.2 1 0
`
		tempCoverageFilePath = "test-profile.cov"
	)

	tempCoverageFile, err := os.OpenFile(tempCoverageFilePath, os.O_WRONLY|os.O_CREATE, 0600)
	assert.NoError(t, err)
	defer func() {
		_ = os.Remove(tempCoverageFilePath)
	}()
	_, err = tempCoverageFile.WriteString(sampleCoverageFile)
	assert.NoError(t, err)
	err = tempCoverageFile.Close()
	assert.NoError(t, err)

	exclusionsMatcher := regexp.MustCompile(`/clean\.go`)
	err = filterCoverage(tempCoverageFilePath, exclusionsMatcher)
	assert.NoError(t, err)

	actualCoverageFile, err := ioutil.ReadFile(tempCoverageFilePath)
	assert.NoError(t, err)
	assert.Equal(t, filteredCoverageFile, string(actualCoverageFile))
}
