package tmpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMethods(t *testing.T) {
	src := `package test

type myType interface {
	LoadByID(ctx context.Context, id int64) myType

`
	srcAST := getASTFromSrc(src)

	result := GetMethods(srcAST, "myType")

	expected1 := []Method{
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
					"myType",
				},
			},
		},
	}

	assert.Equal(t, expected1, result)
}
