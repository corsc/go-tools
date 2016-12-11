package generator

import "regexp"

// UnknownPackage ...
const UnknownPackage = "unknown"

// Coverage will generate coverage for the supplied directory and any sub-directories that contain Go files
func Coverage(basePath string, dirMatcher, fileMatcher *regexp.Regexp, goTestArgs []string) {
	processAllDirs(basePath, dirMatcher, "coverage", func(path string) {
		generateCoverage(path, fileMatcher, goTestArgs)
	})
}

// CoverageSingle will generate coverage for the supplied directory (and ignore all sub directories)
func CoverageSingle(basePath string, fileMatcher *regexp.Regexp, goTestArgs []string) {
	generateCoverage(basePath, fileMatcher, goTestArgs)
}
