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
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

// FileExists returns true if the filename exists and false otherwise.
func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// GetGoFiles returns all the go files from those files supplied
func GetGoFiles(filenames ...string) ([]string, error) {
	files := []string{}

	for _, filename := range filenames {
		if strings.HasSuffix(filename, ".go") {
			files = append(files, filename)
		} else {
			return nil, fmt.Errorf("'%s' is not a Go file", filename)
		}
	}

	return files, nil
}

// GetGoFilesFromCurrentDir returns all the Go files in the current dir
func GetGoFilesFromCurrentDir() ([]string, error) {
	return GetGoFilesFromDir(".")
}

// GetGoFilesFromDir returns all the Go files in the supplied dir
func GetGoFilesFromDir(dirname string) ([]string, error) {
	pkg, err := build.ImportDir(dirname, 0)
	if err != nil {
		return []string{}, err
	}

	return getGoFilesFromPackage(pkg, err)
}

// GetGoFilesFromDirectoryRecursive returns all the Go files from the supplied directory and it's children
func GetGoFilesFromDirectoryRecursive(dirname string) ([]string, error) {
	files := []string{}

	if !strings.HasSuffix(dirname, "...") {
		dirname += "..."
	}

	packages := GetAllPackagesUnderDirectory(dirname)
	if len(packages) == 0 {
		return files, fmt.Errorf("no go files found in dir '%s'", dirname)
	}

	for _, dirname := range packages {
		theseFiles, err := GetGoFilesFromDir(dirname)
		if err != nil {
			return []string{}, err
		}
		files = append(files, theseFiles...)
	}

	return files, nil
}

// returns all the go files in the supplied package
func getGoFilesFromPackage(pkg *build.Package, err error) ([]string, error) {
	files := []string{}

	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return files, nil
		}
		return files, err
	}

	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.TestGoFiles...)
	if pkg.Dir != "." {
		for i, f := range files {
			files[i] = filepath.Join(pkg.Dir, f)
		}
	}

	return files, nil
}
