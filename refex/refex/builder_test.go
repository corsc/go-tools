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

func TestCodeBuilder_build(t *testing.T) {
	scenarios := []struct {
		codeIn      string
		beforeParts []*part
		afterParts  []*part
		expected    string
	}{
		{
			codeIn: "statsd.Count(a, b)",
			beforeParts: []*part{
				{code: `something, err := statsd.Count(`},
				{
					code:  `a, b`,
					isArg: true,
					index: 1,
				},
				{code: `)`},
			},
			afterParts: []*part{
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
			expected: `stats.D.Count(a, b)`,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.codeIn, func(t *testing.T) {
			builder := &codeBuilderImpl{}
			result, resultErr := builder.build(scenario.beforeParts, scenario.afterParts)
			assert.Nil(t, resultErr)
			assert.Equal(t, scenario.expected, result)
		})
	}

}
