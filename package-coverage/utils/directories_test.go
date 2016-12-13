package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAllGoDirs(t *testing.T) {
	dir := strings.TrimSuffix(GetCurrentDir(), "package-coverage/utils/")

	path := "../"
	expected := []string{
		dir + "package-coverage/",
		dir + "package-coverage/generator/",
		dir + "package-coverage/parser/",
		dir + "package-coverage/tests/fixtures/path_matcher/",
		dir + "package-coverage/tests/fixtures/path_matcher/excluded/",
		dir + "package-coverage/tests/fixtures/path_matcher/included/",
		dir + "package-coverage/utils/",
	}

	currentDir := GetCurrentDir()
	results, err := FindAllGoDirs(path)
	assert.Equal(t, currentDir, GetCurrentDir(), "expected current working directory to be restored")

	assert.Nil(t, err)
	assert.Equal(t, expected, results)
}
