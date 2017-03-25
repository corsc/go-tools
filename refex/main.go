package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"go/format"

	"github.com/corsc/go-tools/refex/refex"
	"golang.org/x/tools/imports"
)

func main() {
	before := ""
	after := ""
	displayOnly := false
	skipFormat := false

	flag.StringVar(&before, "b", "", "the code pattern before changes")
	flag.StringVar(&after, "a", "", "the code pattern after changes")
	flag.BoolVar(&displayOnly, "d", false, "display code on stdOut instead of updating any files")
	flag.BoolVar(&skipFormat, "f", false, "skip formatting of generated code")
	flag.Parse()

	fileName := flag.Arg(0)
	sanityCheck(before, after, fileName)

	fileContents, err := getFileContents(fileName)
	if err != nil {
		exitWithError(err)
	}

	result, err := refex.Do(fileContents, before, after)
	if err != nil {
		exitWithError(err)
	}

	// format code
	if !skipFormat {
		codeAsBytes := []byte(result)
		codeAsBytes, err = goFmt(codeAsBytes)
		if err != nil {
			exitWithError(err)
		}

		codeAsBytes, err = goImports(fileName, codeAsBytes)
		if err != nil {
			exitWithError(err)
		}
		result = string(codeAsBytes)
	}

	var writer io.Writer
	if displayOnly {
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			exitWithError(err)
		}

		defer func() {
			if err := writer.(io.Closer).Close(); err != nil {
				exitWithError(err)
			}
		}()

	}
	fmt.Fprint(writer, result)
}

func getFileContents(fileName string) (string, error) {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func sanityCheck(before string, after string, fileName string) {
	if fileName == "" {
		exitWithError(errors.New("please include the file or directory name as the last argument"))
	}

	if before == "" {
		exitWithError(errors.New("before pattern cannot be empty"))
	}

	if after == "" {
		exitWithError(errors.New("after pattern cannot be empty"))
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

func goFmt(codeIn []byte) ([]byte, error) {
	formattedCode, err := format.Source(codeIn)
	if err != nil {
		fmt.Fprintf(os.Stdout, "warning: invalid code generated. Err: %s", err)
		return codeIn, err
	}
	return formattedCode, nil
}

func goImports(fileName string, codeIn []byte) ([]byte, error) {
	formattedCode, err := imports.Process(fileName, codeIn, &imports.Options{})
	if err != nil {
		fmt.Fprintf(os.Stdout, "warning: invalid code generated. Err: %s", err)
		return codeIn, err
	}
	return formattedCode, nil
}
