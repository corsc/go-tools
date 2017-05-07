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

package commons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllPackagesUnderDirectory(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected []string
	}{
		{
			desc:  "happy path",
			input: "./testdata/get-go-files/",
			expected: []string{
				"./testdata/get-go-files",
			},
		},
		{
			desc:  "happy path - recursive",
			input: "./testdata/get-go-files/...",
			expected: []string{
				"./testdata/get-go-files",
				"./testdata/get-go-files/c",
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result := GetAllPackagesUnderDirectory(scenario.input)
			assert.Equal(t, scenario.expected, result, scenario.desc)
		})
	}
}
