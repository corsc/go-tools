package generator

import "github.com/corsc/go-tools/package-coverage/utils"

// UnknownPackage ...
const UnknownPackage = "unknown"

// Coverage will generate coverage for the supplied directory and any sub-directories that contain Go files
func Coverage(basePath string) {
	paths, err := utils.FindAllGoDirs(basePath)
	if err != nil {
		return
	}

	for _, path := range paths {
		generateCoverage(path)
	}
}

// CoverageSingle will generate coverage for the supplied directory (and ignore all sub directories)
func CoverageSingle(basePath string) {
	generateCoverage(basePath)
}
