package tmpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMethods_HappyPath(t *testing.T) {
	src := `package test

type myType interface {
	LoadByID(ctx context.Context, id int64) bool
}
`
	srcAST := getASTFromSrc(src)

	result := GetMethods(srcAST, "myType")

	expected := []Method{
		{
			Name: "LoadByID",
			Params: []MethodField{
				{
					[]string{"ctx"},
					"context.Context",
				},
				{
					[]string{"id"},
					"int64",
				},
			},
			Results: []MethodField{
				{
					[]string{},
					"bool",
				},
			},
		},
	}

	assert.Equal(t, expected, result)
}

func TestGetMethods_None(t *testing.T) {
	src := `package test

type myInterface interface {
	LoadByID(ctx context.Context, id int64) bool
}

type myStruct struct {}

`
	srcAST := getASTFromSrc(src)

	result := GetMethods(srcAST, "myStruct")

	expected := []Method{}

	assert.Equal(t, expected, result)
}
