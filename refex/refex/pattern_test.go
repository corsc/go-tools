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

package refex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPattern_build(t *testing.T) {
	scenarios := []struct {
		transform string
		expected  []*part
		expectErr bool
	}{
		{
			transform: "statsd.Count($1$)",
			expected: []*part{
				{
					code: "statsd.Count(",
				},
				{
					isArg: true,
					index: 1,
				},
				{
					code: ")",
				},
			},
			expectErr: false,
		},
		{
			transform: "stats.D.Count($1$)",
			expected: []*part{
				{
					code: "stats.D.Count(",
				},
				{
					isArg: true,
					index: 1,
				},
				{
					code: ")",
				},
			},
			expectErr: false,
		},
		{
			transform: "statsd.Count($1$, $2$)",
			expected: []*part{
				{
					code: "statsd.Count(",
				},
				{
					isArg: true,
					index: 1,
				},
				{
					code: ", ",
				},
				{
					isArg: true,
					index: 2,
				},
				{
					code: ")",
				},
			},
			expectErr: false,
		},
		{
			transform: "Do(context.Background(), $2$)  // config: $1$",
			expected: []*part{
				{
					code: "Do(context.Background(), ",
				},
				{
					isArg: true,
					index: 2,
				},
				{
					code: ")  // config: ",
				},
				{
					isArg: true,
					index: 1,
				},
				{
					code: "",
				},
			},
			expectErr: false,
		},
		{
			transform: "DoAsync($1$, $2$)",
			expected: []*part{
				{
					code: "DoAsync(",
				},
				{
					isArg: true,
					index: 1,
				},
				{
					code: ", ",
				},
				{
					isArg: true,
					index: 2,
				},
				{
					code: ")",
				},
			},
			expectErr: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.transform, func(t *testing.T) {
			pattern := &patternImpl{}
			result, resultErr := pattern.build(scenario.transform)
			assert.Equal(t, scenario.expected, result)
			assert.Equal(t, scenario.expectErr, resultErr != nil, "error was: %v", resultErr)
		})
	}
}

func TestPattern_regex(t *testing.T) {
	scenarios := []struct {
		transform string
		expected  string
	}{
		{
			transform: `statsd.Count($1$)`,
			expected:  `(statsd.Count\()` + wildcard + `(\))`,
		},
		{
			transform: `stats.D.Count($1$)`,
			expected:  `(stats.D.Count\()` + wildcard + `(\))`,
		},
		{
			transform: `statsd.Count($1$, $2$)`,
			expected:  `(statsd.Count\()` + wildcard + `(, )` + wildcard + `(\))`,
		},
		{
			transform: `Do(context.Background(), $2$)  // config: $1$`,
			expected:  `(Do\(context.Background\(\), )` + wildcard + `(\)  // config: )` + wildcard + `()`,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.transform, func(t *testing.T) {
			pattern := &patternImpl{}
			_, err := pattern.build(scenario.transform)
			assert.Nil(t, err)

			result, resultErr := pattern.regexp()
			assert.Equal(t, scenario.expected, result.String())
			assert.Nil(t, resultErr)
		})
	}
}
