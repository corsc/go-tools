package generator

import (
	"log"
	"os"
	"regexp"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// Clean will search the supplied directory and any sub-directories that contain Go files and remove any
// existing coverage files
func Clean(basePath string, exclusionsMatcher *regexp.Regexp) {
	processAllDirs(basePath, exclusionsMatcher, "clean", clean)
}

// CleanSingle will search the supplied directory that contain Go files and remove any existing coverage file
func CleanSingle(path string) {
	clean(path)
}

func clean(path string) {
	coverageFile := path + coverageFilename
	removeCoverageFile(coverageFile)
}

// remove the previously created coverage file
func removeCoverageFile(filename string) {
	if _, err := os.Stat(filename); err == nil {
		utils.LogWhenVerbose("removing coverage file @ %s", filename)

		err := os.Remove(filename)
		if err != nil {
			log.Printf("error while removing test file @ %s, err: %s", filename, err)
		}
	}
}
