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

func TestIsSlice_TDT(t *testing.T) {
	scenarios := []struct {
		desc     string
		in       Field
		expected bool
	}{
		{
			desc: "Is slice",
			in: Field{
				Name: "Fu",
				Type: "[]Fus",
			},
			expected: true,
		},
		{
			desc: "Is NOT slice",
			in: Field{
				Name: "Bar",
				Type: "string",
			},
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		result := isSlice(scenario.in)
		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}

func TestIsMap_TDT(t *testing.T) {
	scenarios := []struct {
		desc     string
		in       Field
		expected bool
	}{
		{
			desc: "Is map",
			in: Field{
				Name: "Fu",
				Type: "map[string]Fus",
			},
			expected: true,
		},
		{
			desc: "Is NOT map",
			in: Field{
				Name: "Bar",
				Type: "string",
			},
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		result := isMap(scenario.in)
		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}

func TestParamsWithType(t *testing.T) {
	expected := "a, b int, c string"
	method := Method{
		Name: "fubar",
		Params: []MethodField{
			{
				Names: []string{"a", "b"},
				Type:  "int",
			},
			{
				Names: []string{"c"},
				Type:  "string",
			},
		},
	}

	result := paramsWithType(method)
	assert.Equal(t, expected, result)
}

func TestParamsNoType(t *testing.T) {
	expected := "a, b, c"
	method := Method{
		Name: "fubar",
		Params: []MethodField{
			{
				Names: []string{"a", "b"},
				Type:  "int",
			},
			{
				Names: []string{"c"},
				Type:  "string",
			},
		},
	}

	result := paramsNoType(method)
	assert.Equal(t, expected, result)
}
