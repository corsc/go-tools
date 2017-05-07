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
