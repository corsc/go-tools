package generator

import "regexp"

// UnknownPackage ...
const UnknownPackage = "unknown"

// Coverage will generate coverage for the supplied directory and any sub-directories that contain Go files
func Coverage(basePath string, matcher *regexp.Regexp) {
	processAllDirs(basePath, matcher, "coverage", generateCoverage)
}

// CoverageSingle will generate coverage for the supplied directory (and ignore all sub directories)
func CoverageSingle(basePath string) {
	generateCoverage(basePath)
}
