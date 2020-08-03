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
	"io"
	"log"
	"strings"
	"text/template"
)

const (
	invalidStr = "invalid"

	NoopTemplate = `{{- $typeName := firstUpper .TypeName }}
package {{ .PackageName }}

type noop{{ $typeName }} struct {}
{{- range $methodIndex, $method := .Methods }}
{{- $paramsLen := len $method.Params }}
{{- $resultsLen := len $method.Results }}

func (*noop{{ $typeName }}) {{ $method.Name }}(
	{{- range $paramIndex, $param := $method.Params }}
        {{- $paramNamesLen := len $param.Names }}
		{{- range $paramNameIndex, $paramName := $param.Names }}_{{ isNotLast $paramNamesLen $paramNameIndex ", " }}{{ end }} {{ $param.Type }}{{ isNotLast $paramsLen $paramIndex ", " }}
	{{- end -}}
) {{ if ne $resultsLen 0 -}}
(
	{{- range $resultIndex, $result := $method.Results }}
        {{- $resultNamesLen := len $result.Names }}
		{{- if eq $resultNamesLen 0 }}_{{ end }}
		{{- range $resultNameIndex, $resultName := $result.Names }}_{{ isNotLast $resultNamesLen $resultNameIndex ", " }}{{ end }} {{ $result.Type }}{{ isNotLast $resultsLen $resultIndex ", " }}
	{{- end -}}
) {{ end -}}
{ {{- if ne $resultsLen 0 }} return {{ end -}} }
{{- end }}

`
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
		"isNotFirst":             isNotFirst,
		"isNotLast":              isNotLast,
		"firstLower":             firstLower,
		"firstUpper":             strings.Title,
		"toUpper":                strings.ToUpper,
		"toLower":                strings.ToLower,
		"isSlice":                isSlice,
		"sliceType":              sliceType,
		"isMap":                  isMap,
		"add":                    add,
		"paramsWithType":         paramsWithType,
		"paramsNoType":           paramsNoType,
		"hasField":               hasField,
		"testData":               testData,
		"fieldsList":             fieldsList,
		"fieldsListWithTag":      fieldsListWithTag,
		"fieldsListWithTagValue": fieldsListWithTagValue,
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

func processFieldsList(rawTemplate string, fields []Field) (string, error) {
	compiledTemplate, err := getTemplate().Parse(rawTemplate)
	if err != nil {
		return "", err
	}

	vars := TemplateData{
		Fields: fields,
	}

	buffer := &bytes.Buffer{}
	err = compiledTemplate.Execute(buffer, vars)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func fieldsList(fields []Field, fragment string) (string, error) {
	rawTemplate := `
{{- $len := len .Fields }}
	{{- range $index, $field := .Fields -}}
	` + fragment + `{{ isNotLast $len $index ", " }}
{{- end }}`

	return processFieldsList(rawTemplate, fields)
}

func fieldsListWithTag(fields []Field, fragment, tag string) (string, error) {
	rawTemplate := `
{{- $len := len .Fields }}
	{{- range $index, $field := .Fields -}}
	` + fragment + `{{ isNotLast $len $index ", " }}
{{- end }}`

	var filteredFields []Field
	for _, thisField := range fields {
		if _, hasTag := thisField.Tags[tag]; hasTag {
			filteredFields = append(filteredFields, thisField)
		}
	}

	return processFieldsList(rawTemplate, filteredFields)
}

func fieldsListWithTagValue(fields []Field, fragment, tag, value string) (string, error) {
	rawTemplate := `
{{- $len := len .Fields }}
	{{- range $index, $field := .Fields -}}
	` + fragment + `{{ isNotLast $len $index ", " }}
{{- end }}`

	var filteredFields []Field

	for _, thisField := range fields {
		if thisValue, _ := thisField.Tags[tag]; thisValue == value {
			filteredFields = append(filteredFields, thisField)
		}
	}

	return processFieldsList(rawTemplate, filteredFields)
}

// generate predictable test data based on the type and index
func testData(index int, destType string) string {
	switch destType {
	case "int":
		return intTestData[index%10]

	case "int8", "int16", "int32", "int64":
		return destType + "(" + intTestData[index%10] + ")"

	case "float32":
		return destType + "(" + floatTestData[index%10] + ")"

	case "float64":
		return floatTestData[index%10]

	case "string":
		return stringTestData[index%10]

	case "[]byte":
		return byteTestData[index%10]
	}

	switch {
	case strings.HasPrefix(destType, "*"):
		return "nil"

	default:
		return destType + "{}"
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
