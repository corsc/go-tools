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
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLinePosFromPos_happyPath(t *testing.T) {
	scenarios := []struct {
		desc          string
		source        string
		pos           int
		expectedStart int
		expectedEnd   int
	}{
		{
			desc: "test 1",
			source: `package main

func main() {}
`,
			pos:           1,
			expectedStart: 0,
			expectedEnd:   13,
		},
		{
			desc: "line with only a linebreak",
			source: `package main

func main() {}
`,
			pos:           13,
			expectedStart: 13,
			expectedEnd:   13,
		},
		{
			desc: "test 3",
			source: `package main

func main() {}
`,
			pos:           18,
			expectedStart: 14,
			expectedEnd:   29,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			resultStart, resultEnd := GetLineBoundary([]byte(scenario.source), token.Pos(scenario.pos))

			assert.Equal(t, scenario.expectedStart, resultStart)
			assert.Equal(t, scenario.expectedEnd, resultEnd)
		})
	}
}

func TestGetLineBoundary_invalidPos(t *testing.T) {
	scenarios := []struct {
		desc   string
		source string
		pos    int
	}{
		{
			desc:   "empty source",
			source: ``,
			pos:    0,
		},
		{
			desc: "invalid pos",
			source: `package main

func main() {}
`,
			pos: 666,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			assert.Panics(t, func() {
				_, _ = GetLineBoundary([]byte(scenario.source), token.Pos(scenario.pos))
			})
		})
	}
}
