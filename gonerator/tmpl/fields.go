package tmpl

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"
)

// Field ...
type Field struct {
	Name string
	Type string
	Tags map[string]string

	Fields []Field
}

// String implements the stringer interface
func (f Field) String() string {
	return fmt.Sprintf("Name: %s (type: %s)", f.Name, f.Type)
}

// GetFields will extract a slice of fields from the supplied AST
func GetFields(file *ast.File, typeName string) []Field {
	out := []Field{}

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					if spec.Name.Name != typeName {
						continue
					}

					switch sType := spec.Type.(type) {

					case *ast.StructType:
						for _, field := range sType.Fields.List {
							name := field.Names[0].Name
							typeName := getTypeString(field.Type)

							thisField := Field{
								Name: name,
								Type: typeName,
							}

							if field.Tag != nil {
								sTag := structTag(field.Tag.Value)
								thisField.Tags = sTag.getAll()
							}

							// process for custom structs
							subFields := GetFields(file, typeName)
							if len(subFields) > 0 {
								thisField.Fields = subFields
							}

							out = append(out, thisField)
						}

					default:
						continue
					}

				default:
					continue
				}
			}

		default:
			continue
		}
	}

	return out
}

// Copied from the reflect package
//
// A StructTag is the tag string in a struct field.
//
// By convention, tag strings are a concatenation of
// optionally space-separated key:"value" pairs.
// Each key is a non-empty string consisting of non-control
// characters other than space (U+0020 ' '), quote (U+0022 '"'),
// and colon (U+003A ':').  Each value is quoted using U+0022 '"'
// characters and Go string literal syntax.
type structTag string

// Get returns all the tags as a map
func (tag structTag) getAll() map[string]string {
	tag = structTag(strings.Replace(string(tag), "`", "", -1))

	out := make(map[string]string)

	for tag != "" {
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			break
		}
		out[name] = value
	}
	return out
}
