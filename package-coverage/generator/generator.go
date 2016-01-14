package generator

const UnknownPackage = "unknown"

// Coverage will generate coverage for the supplied directory and any sub-directories that contain Go files
func Coverage(basePath string) {
	paths, err := findAllGoDirs(basePath)
	if err != nil {
		return
	}

	for _, path := range paths {
		generateCoverage(path)
	}
}
