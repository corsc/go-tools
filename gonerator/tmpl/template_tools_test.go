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

func TestIsNotLast(t *testing.T) {
	scenarios := []struct {
		desc     string
		len      int
		index    int
		insert   string
		expected string
	}{
		{
			desc:     "empty list",
			len:      0,
			index:    0,
			insert:   "FU",
			expected: "FU",
		},
		{
			desc:     "last",
			len:      3,
			index:    2,
			insert:   "FU",
			expected: "",
		},
		{
			desc:     "not last",
			len:      3,
			index:    1,
			insert:   "FU",
			expected: "FU",
		},
	}

	for _, scenario := range scenarios {
		result := isNotLast(scenario.len, scenario.index, scenario.insert)
		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}

func TestIsNotFirst(t *testing.T) {
	scenarios := []struct {
		desc     string
		len      int
		index    int
		insert   string
		expected string
	}{
		{
			desc:     "empty list",
			len:      0,
			index:    0,
			insert:   "FU",
			expected: "",
		},
		{
			desc:     "first",
			len:      3,
			index:    0,
			insert:   "FU",
			expected: "",
		},
		{
			desc:     "not first",
			len:      3,
			index:    2,
			insert:   "FU",
			expected: "FU",
		},
	}

	for _, scenario := range scenarios {
		result := isNotFirst(scenario.len, scenario.index, scenario.insert)
		assert.Equal(t, scenario.expected, result, scenario.desc)
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
