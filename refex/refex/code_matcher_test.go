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
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeMatcher_find(t *testing.T) {
	scenarios := []struct {
		code     string
		pattern  pattern
		expected []*match
	}{
		{
			code: "something, err := statsd.Count(a, b)",
			pattern: &patternStub{
				regex: regexp.MustCompile(`(statsd.Count\()` + wildcard + `(\))`),
			},
			expected: []*match{
				{
					startPos: 18,
					endPos:   36,
					pattern:  `(statsd.Count\()` + wildcard + `(\))`,
				},
			},
		},
		{
			code: "something, err := stats.D.Count(a, b)",
			pattern: &patternStub{
				regex: regexp.MustCompile(`(stats.D.Count\()` + wildcard + `(\))`),
			},
			expected: []*match{
				{
					startPos: 18,
					endPos:   37,
					pattern:  `(stats.D.Count\()` + wildcard + `(\))`,
				},
			},
		},
		{
			code: "something, err := statsd.Count($1$, $2$)",
			pattern: &patternStub{
				regex: regexp.MustCompile(`(statsd.Count\()` + wildcard + `(, )` + wildcard + `(\))`),
			},
			expected: []*match{
				{
					startPos: 18,
					endPos:   40,
					pattern:  `(statsd.Count\()` + wildcard + `(, )` + wildcard + `(\))`,
				},
			},
		},
		{
			code: `package examples

import (
	"fmt"
	"math/rand"
)

func Example1() {
	fmt.Printf("Roll: %d", rand.Intn(6))
	fmt.Printf("Roll: %d", rand.Intn(10))
	fmt.Printf("Roll: %d", rand.Intn(12))
}
`,
			pattern: &patternStub{
				regex: regexp.MustCompile(`(rand.Intn\()` + wildcard + `(\))`),
			},
			expected: []*match{
				{
					startPos: 92,
					endPos:   104,
					pattern:  `(rand.Intn\()` + wildcard + `(\))`,
				},
				{
					startPos: 130,
					endPos:   143,
					pattern:  `(rand.Intn\()` + wildcard + `(\))`,
				},
				{
					startPos: 169,
					endPos:   182,
					pattern:  `(rand.Intn\()` + wildcard + `(\))`,
				},
			},
		},
		{
			code: `package examples

				func Example1() {
					_, err = DoAsync(masterConfig, "SETEX", key, int64(ttl.Seconds()), raw)
				}
		`,
			pattern: &patternStub{
				regex: regexp.MustCompile(`(DoAsync\()` + wildcard + `(, )` + wildcard + `(\))`),
			},
			expected: []*match{
				{
					startPos: 54,
					endPos:   116,
					pattern:  `(DoAsync\()` + wildcard + `(, )` + wildcard + `(\))`,
				},
			},
		},
		{
			code: `v, err := redis.Bytes(DoAsync(r.SlaveRedis(), "GET", key))`,
			pattern: &patternStub{
				regex: regexp.MustCompile(`(DoAsync\()` + wildcard + `(, )` + wildcard + `(\))`),
			},
			expected: []*match{
				{
					startPos: 22,
					endPos:   57,
					pattern:  `(DoAsync\()` + wildcard + `(, )` + wildcard + `(\))`,
				},
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.code, func(t *testing.T) {
			matcher := &codeMatcherImpl{}

			resultErr := matcher.find(scenario.code, scenario.pattern)
			assert.Nil(t, resultErr)
			assert.Equal(t, scenario.expected, matcher.matches, "first match was: '%s'", scenario.code[matcher.matches[0].startPos:matcher.matches[0].endPos])
		})
	}
}

func TestCodeMatcher_buildParts(t *testing.T) {
	scenarios := []struct {
		codeIn   string
		matches  []*match
		expected []*part
	}{
		{
			codeIn: "something, err := statsd.Count(a, b)",
			matches: []*match{
				{
					startPos: 18,
					endPos:   36,
					pattern:  `(statsd.Count\()` + wildcard + `(\))`,
				},
			},
			expected: []*part{
				{code: `something, err := statsd.Count(`},
				{
					code:  `a, b`,
					isArg: true,
					index: 1,
				},
				{code: `)`},
			},
		},
		{
			codeIn: "something, err := stats.D.Count(a, b)",
			matches: []*match{
				{
					startPos: 18,
					endPos:   37,
					pattern:  `(stats.D.Count\()` + wildcard + `(\))`,
				},
			},
			expected: []*part{
				{code: `something, err := stats.D.Count(`},
				{
					code:  `a, b`,
					isArg: true,
					index: 1,
				},
				{code: `)`},
			},
		},
		{
			codeIn: "something, err := statsd.Count(a, b)",
			matches: []*match{
				{
					startPos: 18,
					endPos:   36,
					pattern:  `(statsd.Count\()` + wildcard + `(, )` + wildcard + `(\))`,
				},
			},
			expected: []*part{
				{code: `something, err := statsd.Count(`},
				{
					code:  `a`,
					isArg: true,
					index: 1,
				},
				{code: `, `},
				{
					code:  `b`,
					isArg: true,
					index: 2,
				},
				{code: `)`},
			},
		},
		{
			codeIn: `package examples

				func Example1() {
					_, err = DoAsync(masterConfig, "SETEX", key, int64(ttl.Seconds()), raw)
				}
				`,
			matches: []*match{
				{
					startPos: 54,
					endPos:   116,
					pattern:  `(DoAsync\()` + wildcard + `(, )` + wildcard + `(\))`,
				},
			},
			expected: []*part{
				{code: `package examples

				func Example1() {
					_, err = DoAsync(`},
				{
					code:  `masterConfig`,
					isArg: true,
					index: 1,
				},
				{code: `, `},
				{
					code:  `"SETEX", key, int64(ttl.Seconds()), raw`,
					isArg: true,
					index: 2,
				},
				{code: `)`},
			},
		},
		{
			codeIn: `v, err := redis.Bytes(DoAsync(r.SlaveRedis(), "GET", key))`,
			matches: []*match{
				{
					startPos: 22,
					endPos:   57,
					pattern:  `(DoAsync\()` + wildcard + `(, )` + wildcard + `(\))`,
				},
			},
			expected: []*part{
				{code: `v, err := redis.Bytes(DoAsync(`},
				{
					code:  `r.SlaveRedis()`,
					isArg: true,
					index: 1,
				},
				{code: `, `},
				{
					code:  `"GET", key`,
					isArg: true,
					index: 2,
				},
				{code: `)`},
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.codeIn, func(t *testing.T) {
			matcher := &codeMatcherImpl{
				code:    scenario.codeIn,
				matches: scenario.matches,
			}

			resultErr := matcher.buildParts()
			assert.Nil(t, resultErr)
			assert.Equal(t, scenario.expected, matcher.matches[0].parts)
		})
	}
}

// implements pattern interface
type patternStub struct {
	regex *regexp.Regexp
}

// implements pattern
func (p *patternStub) build(in string) ([]*part, error) {
	return nil, errors.New("not implemented")
}

// implements pattern
func (p *patternStub) regexp() (*regexp.Regexp, error) {
	return p.regex, nil
}
