package generator

import (
	"log"
	"os"

	"regexp"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// Clean will search the supplied directory and any sub-directories that contain Go files and remove any
// existing coverage files
func Clean(basePath string, matcher *regexp.Regexp) {
	processAllDirs(basePath, matcher, clean)
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
