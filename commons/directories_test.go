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

func TestIsDir(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected bool
	}{
		{
			desc:     "happy path",
			input:    "./testdata/file-exists/",
			expected: true,
		},
		{
			desc:     "doesn't exist",
			input:    "./testdata/doesnt-exist/",
			expected: false,
		},
		{
			desc:     "not a directory",
			input:    "./testdata/file-exists/exists.go",
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result := IsDir(scenario.input)
			assert.Equal(t, scenario.expected, result, scenario.desc)
		})
	}

}
