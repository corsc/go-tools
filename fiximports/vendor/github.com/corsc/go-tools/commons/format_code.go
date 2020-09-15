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

package commons

import (
	"go/format"

	"golang.org/x/tools/imports"
)

// GoFmt will format the supplied code using gofmt
func GoFmt(codeIn []byte) ([]byte, error) {
	formattedCode, err := format.Source(codeIn)
	if err != nil {
		return codeIn, err
	}

	return formattedCode, nil
}

// GoImports will format the supplied code using goimports
func GoImports(fileName string, codeIn []byte) ([]byte, error) {
	options := &imports.Options{
		AllErrors: true,
		Comments:  true,
	}

	formattedCode, err := imports.Process(fileName, codeIn, options)
	if err != nil {
		return codeIn, err
	}

	return formattedCode, nil
}
