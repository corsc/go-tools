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
		dir + "package-coverage/utils/",
	}

	results, err := FindAllGoDirs(path)

	assert.Nil(t, err)
	assert.Equal(t, expected, results)
}
