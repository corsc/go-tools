package generator

import (
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func findAllGoDirs(basePath string) ([]string, error) {
	found := []string{}

	err := os.Chdir(basePath)
	if err != nil {
		return nil, err
	}

	filepath.Walk("./", func(path string, finfo os.FileInfo, err error) error {
		if err != nil {
			log.Printf("failed to check path '%s' with error %s", path, err)
			return nil
		}

		if !finfo.IsDir() {
			return nil
		}

		_, filename := filepath.Split(path)
		if strings.HasPrefix(filename, ".") || strings.HasPrefix(filename, "_") || filename == "testdata" {
			return filepath.SkipDir
		}

		pathEnd := getPathEnd(path)

		if hiddenOrSystemDirs(pathEnd) {
			return filepath.SkipDir
		}

		if hasGoFiles(path) {
			if path == "./" {
				found = append(found, getCurrentDir())
			} else {
				path := getCurrentDir() + path + "/"
				found = append(found, path)
			}
		}

		return nil
	})

	return found, nil
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

func getCurrentDir() string {
	absPath, _ := os.Getwd()
	return absPath + "/"
}
