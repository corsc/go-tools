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

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/corsc/go-tools/commons"
	"github.com/corsc/go-tools/refex/refex"
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
		codeAsBytes, err = commons.GoFmt(codeAsBytes)
		if err != nil {
			exitWithError(err)
		}

		codeAsBytes, err = commons.GoImports(fileName, codeAsBytes)
		if err != nil {
			exitWithError(err)
		}
		result = string(codeAsBytes)
	}

	var writer io.WriteCloser
	if settings.displayOnly {
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			exitWithError(err)
		}

		defer commons.CloseIO(writer)
	}

	_, err = fmt.Fprint(writer, result)
	if err != nil {
		commons.LogError("failed to write to writer with err: %s", err)
	}
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
	commons.LogError("Error: %s\n", err)
	os.Exit(1)
}
