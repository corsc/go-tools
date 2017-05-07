// Copyright 2017 Corey Scott http://www.sage42.org/
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"regexp"
)

// UnknownPackage ...
const UnknownPackage = "unknown"

// Coverage will generate coverage for the supplied directory and any sub-directories that contain Go files
func Coverage(basePath string, exclusionsMatcher *regexp.Regexp, quiet bool, goTestArgs []string) {
	processAllDirs(basePath, exclusionsMatcher, "coverage", func(path string) {
		generateCoverage(path, exclusionsMatcher, quiet, goTestArgs)
	})
}

// CoverageSingle will generate coverage for the supplied directory (and ignore all sub directories)
func CoverageSingle(basePath string, exclusionsMatcher *regexp.Regexp, quiet bool, goTestArgs []string) {
	generateCoverage(basePath, exclusionsMatcher, quiet, goTestArgs)
}
