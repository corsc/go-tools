package tmpl

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

// TemplateData is the data structure passed to the template engine
type TemplateData struct {
	TypeName string

	Fields []Field
}

// Generate will populate the buffer with generated code using the supplied type and template
func Generate(buffer *bytes.Buffer, data TemplateData, templateContent string) {
	masterTmpl, err := getTemplate().Parse(templateContent)
	if err != nil {
		log.Fatalf("eror whil parsing template: %s", err)
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
		"isNotLast":  isNotLast,
		"firstLower": firstLower,
		"firstUpper": strings.Title,
	}
}

func isNotLast(len int, index int, insert string) string {
	if (len - 1) == index {
		return ""
	}
	return insert
}

func getTypeData(typeName string) *TemplateData {
	return &TemplateData{
		TypeName: typeName,
	}
}

func firstLower(typeName string) string {
	return strings.ToLower(typeName[:1]) + typeName[1:]
}
