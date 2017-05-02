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
