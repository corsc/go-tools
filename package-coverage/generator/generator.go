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

	"github.com/corsc/go-tools/package-coverage/utils"
)

// UnknownPackage ...
const UnknownPackage = "unknown"

// GeneratorDo defines the call API of the generators
type GeneratorDo interface {
	Do()
}

// SingleDirGenerator will generate coverage for a single directory (not recursive)
type SingleDirGenerator struct {
	Generator
}

// implements pathBuilder interface
func (g *SingleDirGenerator) Do() {
	g.do([]string{g.BasePath})
}

// RecursiveGenerator will recursively generated coverage for a tree of directories
type RecursiveGenerator struct {
	Generator
}

// implements pathBuilder interface
func (g *RecursiveGenerator) Do() {
	paths := []string{}

	paths, err := utils.FindAllGoDirs(g.BasePath)
	if err != nil {
		return
	}

	for _, path := range paths {
		if g.Exclusion.FindString(path) != "" {
			utils.LogWhenVerbose("[coverage] path '%s' skipped due to skipDir regex '%s'", path, g.Exclusion.String())
			continue
		}

		paths = append(paths, path)
	}

	g.do(paths)
}

// Generator is the basis for other coverage generators
type Generator struct {
	// BasePath directory to generate coverage for
	BasePath string

	// Exclusion is a regexp match allowing you to exclude/skip some directories from coverage calculation.
	// NOTE: this is ignored in single coverage mode
	Exclusion *regexp.Regexp

	// QuietMode controls how verbose the logging is.  Useful for debugging
	QuietMode bool

	// Tags is arguments passed to the go test runner
	Tags string

	// SequentialMode controls if the tests are run in parallel or not.  Original version was sequential
	SequentialMode bool
}

func (g *Generator) do(paths []string) {
	if g.SequentialMode {
		for _, path := range paths {
			generateCoverage(path, g.Exclusion, g.QuietMode, g.Tags)
		}
	} else {
		resultCh := make(chan struct{})
		defer close(resultCh)

		// run in parallel
		for _, path := range paths {
			go func(inPath string, inExclusion *regexp.Regexp, inQuietMode bool, inTags string) {
				generateCoverage(inPath, inExclusion, inQuietMode, inTags)
				resultCh <- struct{}{}
			}(path, g.Exclusion, g.QuietMode, g.Tags)
		}

		// wait until everything is done
		done := 0
		for range resultCh {
			done++
			if done >= len(paths) {
				return
			}
		}
	}
}
