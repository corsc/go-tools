package tmpl

import (
	"bytes"
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
func Generate(buffer *bytes.Buffer, data TemplateData, templateContent string) {
	masterTmpl, err := getTemplate().Parse(templateContent)
	if err != nil {
		log.Fatalf("error while parsing template: %s", err)
	}
	err = masterTmpl.Execute(buffer, data)
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
	}
}

func isNotFirst(len int, index int, insert string) string {
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
