package generator

import "regexp"

// UnknownPackage ...
const UnknownPackage = "unknown"

// Coverage will generate coverage for the supplied directory and any sub-directories that contain Go files
func Coverage(basePath string, dirMatcher, fileMatcher *regexp.Regexp, verbose bool, goTestArgs []string) {
	processAllDirs(basePath, dirMatcher, "coverage", func(path string) {
		generateCoverage(path, fileMatcher, verbose, goTestArgs)
	})
}

// CoverageSingle will generate coverage for the supplied directory (and ignore all sub directories)
func CoverageSingle(basePath string, fileMatcher *regexp.Regexp, verbose bool, goTestArgs []string) {
	generateCoverage(basePath, fileMatcher, verbose, goTestArgs)
}
