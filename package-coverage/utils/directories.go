package utils

import (
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type mode int

const (
	goFiles mode = iota
	coverageFiles
)

// FindAllGoDirs will find all directories below the supplied with go files in them
func FindAllGoDirs(basePath string) ([]string, error) {
	return finder(basePath, goFiles)
}

// FindAllCoverageFiles will find all directories below the supplied one
func FindAllCoverageFiles(basePath string) ([]string, error) {
	return finder(basePath, coverageFiles)
}

func finder(basePath string, searchFor mode) ([]string, error) {
	found := []string{}

	err := os.Chdir(basePath)
	if err != nil {
		return nil, err
	}

	_ = filepath.Walk("./", func(path string, finfo os.FileInfo, err error) error {
		if err != nil {
			log.Printf("failed to check path '%s' with error %s", path, err)
			return nil
		}

		var foundPath string

		switch searchFor {
		case goFiles:
			foundPath, err = checkForGo(path, finfo)

		case coverageFiles:
			foundPath, err = checkForCoverage(path, finfo)
		}

		if err != nil {
			return err
		}
		if foundPath != "" {
			found = append(found, foundPath)
		}
		return nil
	})

	return found, nil
}

func checkForGo(path string, finfo os.FileInfo) (string, error) {
	if !finfo.IsDir() {
		return "", nil
	}

	_, filename := filepath.Split(path)
	if strings.HasPrefix(filename, ".") || strings.HasPrefix(filename, "_") || filename == "testdata" {
		return "", filepath.SkipDir
	}

	pathEnd := getPathEnd(path)

	if hiddenOrSystemDirs(pathEnd) {
		return "", filepath.SkipDir
	}

	if !hasGoFiles(path) {
		return "", nil
	}

	if path == "./" {
		return GetCurrentDir(), nil
	}
	foundPath := GetCurrentDir() + path + "/"
	return foundPath, nil
}

func checkForCoverage(path string, finfo os.FileInfo) (string, error) {
	if finfo.IsDir() {
		return "", nil
	}

	_, filename := filepath.Split(path)
	if strings.HasPrefix(filename, ".") || strings.HasPrefix(filename, "_") {
		return "", nil
	}

	if strings.HasSuffix(path, ".cov") {
		foundPath := GetCurrentDir() + path
		return foundPath, nil
	}
	return "", nil
}

func getPathEnd(path string) string {
	pathPrefix := filepath.Dir(path)
	return strings.TrimPrefix(path, pathPrefix)
}

func hiddenOrSystemDirs(pathEnd string) bool {
	return strings.HasPrefix(pathEnd, "/.") || strings.HasPrefix(pathEnd, "/_")
}

func hasGoFiles(path string) bool {
	if _, err := build.ImportDir(path, 0); err != nil {
		if _, noGo := err.(*build.NoGoError); !noGo {
			log.Print(err)
		}
		return false
	}

	return true
}

// GetCurrentDir will return the currently executing directory
func GetCurrentDir() string {
	absPath, _ := os.Getwd()
	return absPath + "/"
}
