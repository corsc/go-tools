package tmpl

import (
	"bytes"
	"log"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldsAsList(t *testing.T) {
	tmpl := `row.Scan({{$len := len . }}{{range $index, $value := .}}&in.{{$value.Name}}{{isNotLast $len $index ", "}}{{end}})`

	vars := []Field{
		{
			Name: "ID",
			Type: "int64",
		},
		{
			Name: "Name",
			Type: "string",
		},
		{
			Name: "Balance",
			Type: "float64",
		},
	}

	masterTmpl, err := getTemplate().Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}
	buffer := &bytes.Buffer{}
	_ = masterTmpl.Execute(buffer, vars)

	assert.Equal(t, "row.Scan(&in.ID, &in.Name, &in.Balance)", buffer.String())
}

func TestIsNotLastAndIsNotFirst(t *testing.T) {
	scenarios := []struct {
		desc               string
		len                int
		index              int
		insert             string
		expectedIsNotFirst string
		expectedIsNotLast  string
	}{
		{
			desc:               "empty list",
			len:                0,
			index:              0,
			insert:             "BAR",
			expectedIsNotFirst: "BAR",
			expectedIsNotFirst: "",
		},
		{
			desc:               "last",
			len:                6,
			index:              5,
			insert:             "BAR",
			expectedIsNotFirst: "",
			expectedIsNotLast:  "",
		},
		{
			desc:               "not last",
			len:                33,
			index:              22,
			insert:             "FU",
			expectedIsNotFirst: "FU",
			expectedIsNotLast:  "FU",
		},
	}

	for _, scenario := range scenarios {
		resultIsNotLast := isNotLast(scenario.len, scenario.index, scenario.insert)
		assert.Equal(t, scenario.expectedIsNotLast, resultIsNotLast, scenario.desc)

		resultIsNotFirst := isNotFirst(scenario.len, scenario.index, scenario.insert)
		assert.Equal(t, scenario.expectedIsNotFirst, resultIsNotFirst, scenario.desc)
	}
}

func TestFirstLower(t *testing.T) {
	scenarios := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "no change",
			input:    "fu",
			expected: "fu",
		},
		{
			desc:     "ALL CAPS",
			input:    "FU",
			expected: "fU",
		},
		{
			desc:     "all lower",
			input:    "fu",
			expected: "fu",
		},
	}

	for _, scenario := range scenarios {
		result := firstLower(scenario.input)
		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}
