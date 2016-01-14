package generator

import (
	"testing"

	"os"

	"log"

	"github.com/stretchr/testify/assert"
	"sage42.org/go-tools/package-coverage/utils"
)

func init() {
	// change test file so that it doesn't overlap when using this tool on itself
	fakeTestFilename = "my-fake-test.go"
}

func TestGenerateCoverage_HappyPath(t *testing.T) {
	defer restoreExecCoverage(mockExecCoverage())

	path := utils.GetCurrentDir()
	expected := path + "generator.cov"

	result := generateCoverage(path)
	assert.Equal(t, expected, result)
}

func TestAddFakeTest_HappyPath(t *testing.T) {
	path := utils.GetCurrentDir()
	packageName := "generator"
	expectedFilename := path + fakeTestFilename
	defer removeTestFile(expectedFilename)

	result := addFakeTest(path, packageName)
	assert.Equal(t, expectedFilename, result)

	assertFileExists(t, expectedFilename)
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

func TestRemoveFakeTest(t *testing.T) {
	packageName := "mypackage"
	testFile := utils.GetCurrentDir() + "my_test.go"

	createTestFile(packageName, testFile)
	removeFakeTest(testFile)

	_, err := os.OpenFile(testFile, os.O_RDWR, 0666)
	assert.True(t, os.IsNotExist(err))
}

func mockExecCoverage() func(string, string) error {
	original := execCoverage

	execCoverage = func(string, string) error {
		return nil
	}

	return original
}

func restoreExecCoverage(original func(string, string) error) {
	execCoverage = original
}

func removeTestFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Printf("error while removing test file %s", err)
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
