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
			return nil, fmt.Errorf("'%s' is not a Go file")
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
		return nil, err
	}

	return getGoFilesFromPackage(pkg, err)
}

// GetGoFilesFromDirectoryRecursive returns all the Go files from the supplied directory and it's children
func GetGoFilesFromDirectoryRecursive(dirname string) ([]string, error) {
	files := []string{}

	if !strings.HasSuffix(dirname, "...") {
		dirname += "..."
	}

	for _, dirname := range GetAllPackagesUnderDirectory(dirname) {
		theseFiles, err := GetGoFilesFromDir(dirname)
		if err != nil {
			return nil, err
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
		return nil, err
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
