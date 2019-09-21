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

func TestGetMethods(t *testing.T) {
	scenarios := []struct {
		desc     string
		inSrc    string
		inType   string
		expected []Method
	}{
		{
			desc:   "happy path - interface",
			inSrc:  srcSampleInterface,
			inType: "myType",
			expected: []Method{
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
			},
		},
		{
			desc:     "happy path - empty struct",
			inSrc:    srcEmptyStruct,
			inType:   "myStruct",
			expected: nil,
		},
		// Reference: https://github.com/corsc/go-tools/pull/26
		// Thanks @ybdx
		{
			desc:     "happy path - interface with no params or results",
			inSrc:    srcInterfaceWithNoParamsAndNoResults,
			inType:   "myType",
			expected: []Method{
				{
					Name: "MyTest",
					Params: nil,
					Results: nil,
				},
			},
		},
	}

	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.desc, func(t *testing.T) {
			srcAST := getASTFromSrc(scenario.inSrc)

			result := GetMethods(srcAST, scenario.inType)

			assert.Equal(t, scenario.expected, result)
		})
	}
}

var (
	srcSampleInterface = `package test

type myType interface {
	LoadByID(ctx context.Context, id int64) bool
}
`
	srcEmptyStruct = `package test

type myStruct struct {}
`

	srcInterfaceWithNoParamsAndNoResults = `package test

type myType interface { MyTest() }
`
)
