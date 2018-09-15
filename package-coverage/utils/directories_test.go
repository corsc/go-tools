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
		dir + "package-coverage/config/",
		dir + "package-coverage/generator/",
		dir + "package-coverage/parser/",
		dir + "package-coverage/test-data/pathmatcher/",
		dir + "package-coverage/test-data/pathmatcher/excluded/",
		dir + "package-coverage/test-data/pathmatcher/included/",
		dir + "package-coverage/utils/",
	}

	currentDir := GetCurrentDir()
	results, err := FindAllGoDirs(path)
	assert.Equal(t, currentDir, GetCurrentDir(), "expected current working directory to be restored")

	assert.Nil(t, err)
	assert.Equal(t, expected, results)
}
