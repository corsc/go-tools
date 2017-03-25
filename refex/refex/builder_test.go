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
