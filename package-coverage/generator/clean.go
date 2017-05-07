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
	"os"
	"regexp"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// NewCleaner returns an instance of the default Cleaner implmentation
func NewCleaner() Cleaner {
	return &cleanerImpl{
		fsWrapper: &fsWrapperImpl{},
	}
}

// Cleaner will remove any previously generated coverage files
type Cleaner interface {
	// Clean a single directory
	Single(path string)

	// Recursive will clean a directory and all child directories (excluding any matched be the regex)
	Recursive(path string, exclusions *regexp.Regexp)
}

// default implementation of the Cleaner interface
type cleanerImpl struct {
	fsWrapper fsWrapper
}

// Single implements the Cleaner interface
func (cleaner *cleanerImpl) Single(path string) {
	cleaner.clean(path)
}

// Clean will search the supplied directory and any sub-directories that contain Go files and remove any
// existing coverage files
func (cleaner *cleanerImpl) Recursive(path string, exclusions *regexp.Regexp) {
	processAllDirs(path, exclusions, "clean", cleaner.clean)
}

func (cleaner *cleanerImpl) clean(path string) {
	coverageFile := path + coverageFilename

	if cleaner.fsWrapper.Exists(coverageFile) {
		utils.LogWhenVerbose("[cleaner] removing coverage file @ %s", coverageFile)
		cleaner.fsWrapper.Delete(coverageFile)
	}
}

type fsWrapper interface {
	// Returns true if a file exists and false otherwise
	Exists(filename string) bool

	// Remove a file
	Delete(filename string)
}

// default implementation of fsWrapper
type fsWrapperImpl struct{}

// Exists implements fsWrapper
func (fs *fsWrapperImpl) Exists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

// Delete implements fsWrapper
func (fs *fsWrapperImpl) Delete(filename string) {
	err := os.Remove(filename)
	if err != nil {
		utils.LogWhenVerbose("[cleaner] failed to remove %s with err: %s", filename, err)
	}
}
