package tmpl

import (
	"fmt"
	"go/ast"
)

// MethodField defines the input params / output results of a method
type MethodField struct {
	Names []string
	Type  string
}

// Method defines the methods defined in an interface
type Method struct {
	Name    string
	Params  []MethodField
	Results []MethodField
}

// String implements the stringer interface
func (f Method) String() string {
	return fmt.Sprintf("Name: %s (params: %s) (results: %s)", f.Name, f.Params, f.Results)
}

// GetMethods will extract a slice of funcs from the supplied AST
func GetMethods(file *ast.File, typename string) []Method {
	out := []Method{}

	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					sType, ok := spec.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}
					for _, field := range sType.Methods.List {
						switch fnType := field.Type.(type) {
						case *ast.FuncType:
							params, results := extractParamsAndResults(fnType)
							fn := Method{
								Name:    field.Names[0].Name,
								Params:  params,
								Results: results,
							}
							out = append(out, fn)
						}
					}
				}
			}
		}
	}

	return out
}

func extractParamsAndResults(funcType *ast.FuncType) ([]MethodField, []MethodField) {
	params := fieldListToMethodFields(funcType.Params)
	results := fieldListToMethodFields(funcType.Results.List)

	return params, results
}

func fieldListToMethodFields(fieldList *ast.FieldList) []MethodField {
	params := []MethodField{}

	for _, item := range fieldList.List {
		typeStr := getTypeString(item.Type)
		funcField := MethodField{
			Type:  typeStr,
			Names: make([]string, len(item.Names)),
		}

		for i := 0; i < len(item.Names); i++ {
			funcField.Names[i] = item.Names[i].Name
		}

		params = append(params, funcField)
	}

	return params
}

func getTypeString(expr ast.Expr) string {
	var result string
	switch etype := expr.(type) {
	case *ast.SelectorExpr:
		result = fmt.Sprintf("%s.%s", etype.X, etype.Sel)
	default:
		result = fmt.Sprintf("%s", etype)
	}
	return result
}
