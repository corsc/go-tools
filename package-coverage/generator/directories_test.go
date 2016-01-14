package generator

import (
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

func TestFindAllGoDirs(t *testing.T) {
	dir := strings.TrimSuffix(getCurrentDir(), "package-coverage/generator/")

	path := "../../"
	expected := []string{
		dir + "package-coverage/",
		dir + "package-coverage/generator/",
		dir + "package-coverage/parser/",
		dir + "package-coverage/utils/",
	}

	results, err := findAllGoDirs(path)

	assert.Nil(t, err)
	assert.Equal(t, expected, results)
}
