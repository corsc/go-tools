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

package fiximports

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUAT(t *testing.T) {
	scenarios := []struct {
		desc           string
		files          []string
		resultFilename string
	}{
		{
			desc:           "single file",
			files:          []string{"./test-data/test-a.in"},
			resultFilename: "./test-data/test-a.out",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			outputCapture := io.Writer(buffer)

			ProcessFiles(scenario.files, outputCapture)

			expected, err := ioutil.ReadFile(scenario.resultFilename)
			assert.Nil(t, err, scenario.desc)
			assert.Equal(t, expected, buffer.Bytes(), scenario.desc)
		})
	}
}
