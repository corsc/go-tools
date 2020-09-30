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
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFieldsWithSubType(t *testing.T) {
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

type TypeC struct {
	A []string
	B map[string]string
}
`
	srcAST := getASTFromSrc(src)

	result := GetFields(srcAST, "myType")

	expected1 := Field{Name: "ID", Type: "int64", NonArrayType: "int64"}
	expected2 := Field{Name: "Name", Type: "string", NonArrayType: "string"}
	expected3 := Field{Name: "Balance", Type: "float64", NonArrayType: "float64"}
	expected4 := Field{Name: "Tag", Type: "string", NonArrayType: "string", Tags: map[string]string{"outputAs": "fu", "outputUsing": "bar"}}
	expected5 := Field{Name: "Child", Type: "TypeB", NonArrayType: "TypeB", Fields: []Field{
		{Name: "A", Type: "string", NonArrayType: "string"},
		{Name: "B", Type: "string", NonArrayType: "string"},
		{Name: "C", Type: "string", NonArrayType: "string"},
	}}

	assert.Equal(t, 5, len(result))
	assert.Equal(t, expected1, result[0])
	assert.Equal(t, expected2, result[1])
	assert.Equal(t, expected3, result[2])
	assert.Equal(t, expected4, result[3])
	assert.Equal(t, expected5, result[4])
}

func TestGetFieldsWithExoticTypes(t *testing.T) {
	src := `package test

type tinyStruct struct {
	num int
	age string
}

type myType struct {
	A []string
	B map[string]string
	C <-chan map[string]string
	D chan<- map[string]string
	E chan struct{}
	F chan tinyStruct
}
`
	srcAST := getASTFromSrc(src)

	result := GetFields(srcAST, "myType")
	expected1 := Field{Name: "A", Type: "[]string", NonArrayType: "string"}
	expected2 := Field{Name: "B", Type: "map[string]string", NonArrayType: "map[string]string"}
	expected3 := Field{Name: "C", Type: "<-chan map[string]string", NonArrayType: "<-chan map[string]string"}
	expected4 := Field{Name: "D", Type: "chan<- map[string]string", NonArrayType: "chan<- map[string]string"}
	expected5 := Field{Name: "E", Type: "chan struct{}", NonArrayType: "chan struct{}"}
	expected6 := Field{Name: "F", Type: "chan tinyStruct", NonArrayType: "chan tinyStruct"}

	assert.Equal(t, 6, len(result))
	assert.Equal(t, expected1, result[0])
	assert.Equal(t, expected2, result[1])
	assert.Equal(t, expected3, result[2])
	assert.Equal(t, expected4, result[3])
	assert.Equal(t, expected5, result[4])
	assert.Equal(t, expected6, result[5])
}

func getASTFromSrc(src string) *ast.File {
	fs := token.NewFileSet()
	srcAST, _ := parser.ParseFile(fs, "", src, 0)
	return srcAST
}
