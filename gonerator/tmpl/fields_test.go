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
}
`
	srcAST := getASTFromSrc(src)

	result := GetFields(srcAST, "myType")

	expected1 := Field{Name: "ID", Type: "int64"}
	expected2 := Field{Name: "Name", Type: "string"}
	expected3 := Field{Name: "Balance", Type: "float64"}

	assert.Equal(t, 3, len(result))
	assert.Equal(t, expected1, result[0])
	assert.Equal(t, expected2, result[1])
	assert.Equal(t, expected3, result[2])
}

func getASTFromSrc(src string) *ast.File {
	fs := token.NewFileSet()
	srcAST, _ := parser.ParseFile(fs, "", src, 0)
	return srcAST
}
