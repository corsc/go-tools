package main

import (
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"path/filepath"

	"github.com/corsc/go-tools/refex/refex"
	"golang.org/x/tools/imports"
)

type settings struct {
	before      string
	after       string
	displayOnly bool
	skipFormat  bool
}

func main() {
	settings := &settings{}

	flag.StringVar(&settings.before, "b", "", "the code pattern before changes")
	flag.StringVar(&settings.after, "a", "", "the code pattern after changes")
	flag.BoolVar(&settings.displayOnly, "d", false, "display code on stdOut instead of updating any files")
	flag.BoolVar(&settings.skipFormat, "f", false, "skip formatting of generated code")
	flag.Parse()

	pathSupplied := flag.Arg(0)
	sanityCheck(settings, pathSupplied)

	pathInfo, err := os.Stat(pathSupplied)
	if err != nil {
		exitWithError(err)
	}

	paths := []string{}

	if pathInfo.IsDir() {
		files, err := ioutil.ReadDir(pathSupplied)
		if err != nil {
			exitWithError(err)
		}

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".go") {
				paths = append(paths, filepath.Join(pathSupplied, file.Name()))
			}
		}
	} else {
		paths = append(paths, pathSupplied)
	}

	for _, thisPath := range paths {
		do(thisPath, settings)
	}
}

func do(fileName string, settings *settings) {
	result, err := refex.DoFile(fileName, settings.before, settings.after)
	if err != nil {
		exitWithError(err)
	}

	// format code
	if !settings.skipFormat {
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
	if settings.displayOnly {
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

func sanityCheck(settings *settings, fileName string) {
	if fileName == "" {
		exitWithError(errors.New("please include the file or directory name as the last argument"))
	}

	if settings.before == "" {
		exitWithError(errors.New("before pattern cannot be empty"))
	}

	if settings.after == "" {
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
	formattedCode, err := imports.Process(fileName, codeIn, &imports.Options{
		AllErrors: true,
		Comments:  true,
	})
	if err != nil {
		fmt.Fprintf(os.Stdout, "warning: invalid code generated. Err: %s", err)
		return codeIn, err
	}
	return formattedCode, nil
}
