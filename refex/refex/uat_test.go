package refex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUAT(t *testing.T) {
	scenarios := []struct {
		desc     string
		code     string
		before   string
		after    string
		expected string
	}{
		{
			desc: "1 wildcard, 1 param",
			code: `package mypackage

func something() {
	statsd.Count1("call")
	statsd.Count1("me")
	statsd.Count1("baby")
}`,
			before: `statsd.Count1($1$)`,
			after:  `stats.D.Count1($1$)`,
			expected: `package mypackage

func something() {
	stats.D.Count1("call")
	stats.D.Count1("me")
	stats.D.Count1("baby")
}`,
		},
		{
			desc: "2 wildcards, 2 params, same ordering",
			code: `package mypackage

func something() {
	statsd.Count1("don't", "call")
	statsd.Count1("me", "baby")
}`,
			before: `statsd.Count1($1$, $2$)`,
			after:  `stats.D.Count1($1$, $2$)`,
			expected: `package mypackage

func something() {
	stats.D.Count1("don't", "call")
	stats.D.Count1("me", "baby")
}`,
		},
		{
			desc: "2 wildcards, 2 params, changed ordering",
			code: `package mypackage

func something() {
	statsd.Count1("don't", "call")
	statsd.Count1("me", "baby")
}`,
			before: `statsd.Count1($1$, $2$)`,
			after:  `stats.D.Count1($2$, $1$)`,
			expected: `package mypackage

func something() {
	stats.D.Count1("call", "don't")
	stats.D.Count1("baby", "me")
}`,
		},
		{
			desc: "2 wildcards, 2 params, drop 1 param",
			code: `package mypackage

func something() {
	statsd.Count1("don't", "call")
	statsd.Count1("me", "baby")
}`,
			before: `statsd.Count1($1$, $2$)`,
			after:  `stats.D.Count1($2$)`,
			expected: `package mypackage

func something() {
	stats.D.Count1("call")
	stats.D.Count1("baby")
}`,
		},
		{
			desc: "2 wildcards, 2 params, add a param",
			code: `package mypackage

func something() {
	statsd.Count1("don't", "call")
	statsd.Count1("me", "baby")
}`,
			before: `statsd.Count1($1$, $2$)`,
			after:  `stats.D.Count1($2$, "apples", $1$)`,
			expected: `package mypackage

func something() {
	stats.D.Count1("call", "apples", "don't")
	stats.D.Count1("baby", "apples", "me")
}`,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.desc, func(t *testing.T) {
			result, resultErr := Do(scenario.code, scenario.before, scenario.after)
			assert.Nil(t, resultErr)
			assert.Equal(t, scenario.expected, result)
		})
	}
}
