package tmpl

import (
	"go/parser"
	"go/token"
	"testing"

	"go/ast"

	"github.com/stretchr/testify/assert"
)

func TestGetFields(t *testing.T) {
	src := `package test

type myType struct {
	ID      int64
	Name    string
	Balance float64
	Tag string	` + "`" + `outputAs:"fu" outputUsing:"bar"` + "`" + `
	Child TypeB
}

type TypeB struct {
	A string
	B string
	C string
}
`
	srcAST := getASTFromSrc(src)

	result := GetFields(srcAST, "myType")

	expected1 := Field{Name: "ID", Type: "int64"}
	expected2 := Field{Name: "Name", Type: "string"}
	expected3 := Field{Name: "Balance", Type: "float64"}
	expected4 := Field{Name: "Tag", Type: "string", Tags: map[string]string{"outputAs": "fu", "outputUsing": "bar"}}
	expected5 := Field{Name: "Child", Type: "TypeB", Fields: []Field{
		{Name: "A", Type: "string"},
		{Name: "B", Type: "string"},
		{Name: "C", Type: "string"},
	}}

	assert.Equal(t, 5, len(result))
	assert.Equal(t, expected1, result[0])
	assert.Equal(t, expected2, result[1])
	assert.Equal(t, expected3, result[2])
	assert.Equal(t, expected4, result[3])
	assert.Equal(t, expected5, result[4])
}

func getASTFromSrc(src string) *ast.File {
	fs := token.NewFileSet()
	srcAST, _ := parser.ParseFile(fs, "", src, 0)
	return srcAST
}
