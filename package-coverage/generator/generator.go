package generator

import "regexp"

// UnknownPackage ...
const UnknownPackage = "unknown"

// Coverage will generate coverage for the supplied directory and any sub-directories that contain Go files
func Coverage(basePath string, matcher *regexp.Regexp, verbose bool, goTestArgs []string) {
	processAllDirs(basePath, matcher, "coverage", func(path string) { generateCoverage(path, verbose, goTestArgs) })
}

// CoverageSingle will generate coverage for the supplied directory (and ignore all sub directories)
func CoverageSingle(basePath string, verbose bool, goTestArgs []string) {
	generateCoverage(basePath, verbose, goTestArgs)
}
