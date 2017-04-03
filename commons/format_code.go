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
