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
	"strings"
	"sync"

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
	dedupeMap := map[string]struct{}{}

	foundPaths, err := utils.FindAllGoDirs(g.BasePath)
	if err != nil {
		return
	}

	for _, path := range foundPaths {
		if g.Exclusion.FindString(path) != "" {
			utils.LogWhenVerbose("[coverage] path '%s' skipped due to skipDir regex '%s'", path, g.Exclusion.String())
			continue
		}

		// always skip vendor dirs
		if strings.Contains(path, "/vendor/") {
			utils.LogWhenVerbose("[coverage] path '%s' skipped due to /vendor/", path)
			continue
		}

		_, found := dedupeMap[path]
		if !found {
			paths = append(paths, path)
			dedupeMap[path] = struct{}{}
		}
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

	// Race is used to enable --race flag
	Race bool

	// Tags is arguments passed to the go test runner
	Tags string

	// Concurrency controls how many tests can be run concurrently.  Default is `runtime.NumCPU()`
	Concurrency int
}

func (g *Generator) do(paths []string) {
	jobsCh := make(chan string, len(paths))
	wg := &sync.WaitGroup{}

	// create workers
	wg.Add(g.Concurrency)
	for index := 1; index <= g.Concurrency; index++ {
		go doWorker(jobsCh, wg, g.Exclusion, g.QuietMode, g.Race, g.Tags)
	}

	// send the paths
	for _, path := range paths {
		jobsCh <- path
	}
	close(jobsCh)

	for _, path := range paths {
		packageName := findPackageName(path)
		if packageName != UnknownPackage {
			removeFakes(path)
		}
	}

	// wait until everything is done
	wg.Wait()
}

func doWorker(jobsCh <-chan string, wg *sync.WaitGroup, exclusion *regexp.Regexp, quietMode, race bool, tags string) {
	defer wg.Done()

	for path := range jobsCh {
		generateCoverage(path, exclusion, quietMode, race, tags)
	}
}
