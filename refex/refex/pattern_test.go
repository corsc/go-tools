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
	}

	for _, scenario := range scenarios {
		t.Run(scenario.transform, func(t *testing.T) {
			pattern := &patternImpl{}
			pattern.build(scenario.transform)

			result, resultErr := pattern.regexp()
			assert.Equal(t, scenario.expected, result.String())
			assert.Nil(t, resultErr)
		})
	}
}
