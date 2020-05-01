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
	"io"
	"log"
	"strings"
	"text/template"
)

const (
	invalidStr = "invalid"
)

// TemplateData is the data structure passed to the template engine
type TemplateData struct {
	TypeName     string
	PackageName  string
	TemplateFile string
	OutputFile   string

	Fields  []Field
	Extras  []string
	Methods []Method
}

// Generate will populate the buffer with generated code using the supplied type and template
func Generate(writer io.Writer, data TemplateData, templateContent string) {
	masterTmpl, err := getTemplate().Parse(templateContent)
	if err != nil {
		log.Fatalf("error while parsing template: %s", err)
	}
	err = masterTmpl.Execute(writer, data)
	if err != nil {
		log.Fatalf("error while executing template: %s", err)
	}
}

func getTemplate() *template.Template {
	return template.New("master").Funcs(getFuncMap())
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"isNotFirst":     isNotFirst,
		"isNotLast":      isNotLast,
		"firstLower":     firstLower,
		"firstUpper":     strings.Title,
		"toUpper":        strings.ToUpper,
		"toLower":        strings.ToLower,
		"isSlice":        isSlice,
		"sliceType":      sliceType,
		"isMap":          isMap,
		"add":            add,
		"paramsWithType": paramsWithType,
		"paramsNoType":   paramsNoType,
		"hasField":       hasField,
		"testData":       testData,
	}
}

func isNotFirst(_ int, index int, insert string) string {
	if 0 == index {
		return ""
	}
	return insert
}

func isNotLast(len int, index int, insert string) string {
	if (len - 1) == index {
		return ""
	}
	return insert
}

func firstLower(typeName string) string {
	return strings.ToLower(typeName[:1]) + typeName[1:]
}

func isSlice(field Field) bool {
	return strings.HasPrefix(field.Type, "[]")
}

func sliceType(field Field) string {
	if !isSlice(field) {
		return invalidStr
	}

	return strings.Replace(field.Type, "[]", "", 1)
}

func isMap(field Field) bool {
	return strings.Contains(field.Type, "map")
}

func add(inputs ...interface{}) int64 {
	var output int64
	for _, value := range inputs {
		switch concreteValue := value.(type) {
		case int:
			output += int64(concreteValue)

		case int64:
			output += concreteValue

		case int32:
			output += int64(concreteValue)
		}
	}
	return output
}

func paramsWithType(method Method) string {
	return paramsList(method, true)
}

func paramsNoType(method Method) string {
	return paramsList(method, false)
}

func paramsList(method Method, includeType bool) string {
	output := ""

	for outerIndex, param := range method.Params {
		for innerIndex, name := range param.Names {
			output += name + isNotLast(len(param.Names), innerIndex, ", ")
		}
		if includeType {
			output += " " + param.Type
		}
		output += isNotLast(len(method.Params), outerIndex, ", ")
	}

	return output
}

func hasField(fields []Field, fieldName string) bool {
	for _, thisField := range fields {
		if thisField.Name == fieldName {
			return true
		}
	}

	return false
}

// generate predictable test data based on the type and index
func testData(index int, destType string) string {
	switch destType {
	case "int", "int8, int16", "int32", "int64":
		return intTestData[index%10]

	case "float32", "float64":
		return floatTestData[index%10]

	case "string":
		return stringTestData[index%10]

	case "[]byte]":
		return byteTestData[index%10]

	case "time.Time":
		return "time.Time{}"

	default:
		return "nil"
	}
}

var (
	intTestData = []string{
		"111", "222", "333", "444", "555",
		"666", "777", "888", "999", "000",
	}

	floatTestData = []string{
		"1.11", "2.22", "3.33", "4.44", "5.55",
		"6.66", "7.77", "8.88", "9.99", "0.00",
	}

	stringTestData = []string{
		`"AAA"`, `"BBB"`, `"CCC"`, `"DDD"`, `"EEE"`,
		`"FFF"`, `"GGG"`, `"HHH"`, `"III"`, `"JJJ"`,
	}

	byteTestData = []string{
		`[]byte("AAA")`, `[]byte("BBB")`, `[]byte("CCC")`, `[]byte("DDD")`, `[]byte("EEE")`,
		`[]byte("FFF")`, `[]byte("GGG")`, `[]byte("HHH")`, `[]byte("III")`, `[]byte("JJJ")`,
	}
)
