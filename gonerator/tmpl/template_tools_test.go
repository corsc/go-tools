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

func TestIsOrIsNot(t *testing.T) {
	scenarios := []struct {
		desc          string
		funcUnderTest func(len int, index int, insert string) string
		len           int
		index         int
		insert        string
		expected      string
	}{
		{
			desc:          "isNotLast - empty list",
			funcUnderTest: isNotLast,
			len:           0,
			index:         0,
			insert:        "FU",
			expected:      "FU",
		},
		{
			desc:          "isNotLast - last",
			funcUnderTest: isNotLast,
			len:           3,
			index:         2,
			insert:        "FU",
			expected:      "",
		},
		{
			desc:          "isNotLast - not last",
			funcUnderTest: isNotLast,
			len:           3,
			index:         1,
			insert:        "FU",
			expected:      "FU",
		},
		{
			desc:          "isNotFirst - empty list",
			funcUnderTest: isNotFirst,
			len:           0,
			index:         0,
			insert:        "FU",
			expected:      "",
		},
		{
			desc:          "isNotFirst - first",
			funcUnderTest: isNotFirst,
			len:           3,
			index:         0,
			insert:        "FU",
			expected:      "",
		},
		{
			desc:          "isNotFirst - not first",
			funcUnderTest: isNotFirst,
			len:           3,
			index:         2,
			insert:        "FU",
			expected:      "FU",
		},
	}

	for _, scenario := range scenarios {
		result := scenario.funcUnderTest(scenario.len, scenario.index, scenario.insert)
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

func TestIs(t *testing.T) {
	scenarios := []struct {
		desc          string
		funcUnderTest func(Field) bool
		in            Field
		expected      bool
	}{
		{
			desc:          "Is slice",
			funcUnderTest: isSlice,
			in: Field{
				Name: "Fu",
				Type: "[]Fus",
			},
			expected: true,
		},
		{
			desc:          "Is NOT slice",
			funcUnderTest: isSlice,
			in: Field{
				Name: "Bar",
				Type: "string",
			},
			expected: false,
		},
		{
			desc:          "Is map",
			funcUnderTest: isMap,
			in: Field{
				Name: "Fu",
				Type: "map[string]Fus",
			},
			expected: true,
		},
		{
			desc:          "Is NOT map",
			funcUnderTest: isMap,
			in: Field{
				Name: "Bar",
				Type: "string",
			},
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		result := scenario.funcUnderTest(scenario.in)
		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}

func TestParams(t *testing.T) {
	scenarios := []struct {
		desc          string
		expected      string
		method        Method
		funcUnderTest func(method Method) string
	}{
		{
			desc:     "with type",
			expected: "a, b int, c string",
			method: Method{
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
			},
			funcUnderTest: paramsWithType,
		},
		{
			desc:     "with no type",
			expected: "a, b, c",
			method: Method{
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
			},
			funcUnderTest: paramsNoType,
		},
	}

	for _, scenario := range scenarios {
		result := scenario.funcUnderTest(scenario.method)
		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}

func TestTestData(t *testing.T) {
	testDataIndex := 0

	scenarios := []struct {
		desc       string
		inDestType string
		expected   string
	}{
		{
			desc:       "happy path - int",
			inDestType: "int",
			expected:   intTestData[testDataIndex],
		},
		{
			desc:       "happy path - int64",
			inDestType: "int64",
			expected:   "int64(" + intTestData[testDataIndex] + ")",
		},
		{
			desc:       "happy path - float32",
			inDestType: "float32",
			expected:   "float32(" + floatTestData[testDataIndex] + ")",
		},
		{
			desc:       "happy path - float64",
			inDestType: "float64",
			expected:   floatTestData[testDataIndex],
		},
		{
			desc:       "happy path - string",
			inDestType: "string",
			expected:   stringTestData[testDataIndex],
		},
		{
			desc:       "happy path - byte",
			inDestType: "[]byte",
			expected:   byteTestData[testDataIndex],
		},
		{
			desc:       "happy path - time.Time",
			inDestType: "time.Time",
			expected:   "time.Time{}",
		},
		{
			desc:       "happy path - random object",
			inDestType: "foo.Bar",
			expected:   "foo.Bar{}",
		},
		{
			desc:       "happy path - random pointer",
			inDestType: "*foo.Bar",
			expected:   "nil",
		},
	}

	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			result := testData(testDataIndex, scenario.inDestType)
			assert.Equal(t, scenario.expected, result, "expected result")
		})
	}
}

func TestTypeName(t *testing.T) {
	scenarios := []struct {
		desc     string
		in       string
		expected string
	}{
		{
			desc:     "int",
			in:       "int",
			expected: "Int",
		},
		{
			desc:     "int64",
			in:       "int64",
			expected: "Int64",
		},
		{
			desc:     "Bytes",
			in:       "[]byte",
			expected: "SliceOfBytes",
		},
		{
			desc:     "Time",
			in:       "time.Time",
			expected: "TimeTime",
		},
		{
			desc:     "Slice Of String",
			in:       "[]string",
			expected: "SliceOfStrings",
		},
	}

	for _, scenario := range scenarios {
		result := typeName(scenario.in)
		assert.Equal(t, scenario.expected, result, scenario.desc)
	}
}
