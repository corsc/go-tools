package tmpl

import (
	"bytes"
	"log"
	"strings"
	"text/template"
)

// TemplateData is the data structure passed to the template engine
type TemplateData struct {
	TypeName    string
	PackageName string

	Fields  []Field
	Extras  []string
	Methods []Method
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
		"isNotFirst": isNotFirst,
		"isNotLast":  isNotLast,
		"firstLower": firstLower,
		"firstUpper": strings.Title,
		"pbEncode":   pbEncode,
		"pbDecode":   pbDecode,
		"toUpper":    strings.ToUpper,
		"toLower":    strings.ToLower,
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

// Decode helper for Protobuf
func pbDecode(field Field, prefix string) string {
	fieldName := field.Name
	if custom, found := field.Tags["pbName"]; found {
		fieldName = custom
	}

	// Special processing for things like Enum types
	switch field.Type {
	case "time.Time":
		return "time.Unix(0, " + prefix + fieldName + ")"

	default:
		return prefix + fieldName
	}
}

// Encode helper for Protobuf
func pbEncode(field Field, prefix string) string {
	// Special processing for things like Enum types
	switch field.Type {
	case "time.Time":
		return prefix + field.Name + ".UnixNano()"

	default:
		return prefix + field.Name
	}
}
