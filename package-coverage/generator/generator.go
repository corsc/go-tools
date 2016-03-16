package generator

import (
	"log"
	"regexp"
)

// UnknownPackage ...
const UnknownPackage = "unknown"

// Coverage will generate coverage for the supplied directory and any sub-directories that contain Go files
func Coverage(basePath string, matcher *regexp.Regexp) {
	processAllDirs(basePath, matcher, generateCoverage)
}

// CoverageSingle will generate coverage for the supplied directory (and ignore all sub directories)
func CoverageSingle(basePath string, matcher *regexp.Regexp) {
	if matcher.FindString(basePath) != "" {
		log.Printf("Generation of coverage for path '%s' skipped due to skipDir regex '%s'", basePath, matcher.String())
		return
	}

	generateCoverage(basePath)
}
