package tmpl

import (
	"fmt"
	"go/ast"
)

// Field ...
type Field struct {
	Name string
	Type string
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
							typeName := fmt.Sprintf("%s", field.Type)
							out = append(out, Field{Name: name, Type: typeName})
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
