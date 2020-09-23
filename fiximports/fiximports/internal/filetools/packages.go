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

package filetools

/*
This file holds a copy of the import path matching code of
https://code.google.com/p/go/source/browse/src/cmd/go/main.go. It can be
replaced when https://code.google.com/p/go/issues/detail?id=8768 is resolved.
*/

import (
	"github.com/corsc/go-tools/fiximports/fiximports/internal/log"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// GetAllPackagesUnderDirectory is like allPackages but is passed a pattern
// beginning ./ or ../, meaning it should scan the tree rooted
// at the given directory.  There are sometimes ... in the pattern too.
func GetAllPackagesUnderDirectory(pattern string) []string {
	pkgs := matchPackagesInFS(pattern)
	if len(pkgs) == 0 {
		log.Error("warning: %q matched no packages\n", pattern)
	}
	return pkgs
}

func matchPackagesInFS(pattern string) []string {
	dir := getStartDir(pattern)

	prefix, match := buildMatcher(pattern)

	var pkgs []string
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if err != nil || !fi.IsDir() {
			return nil
		}
		if path == dir {
			// filepath.Walk starts at dir and recurses. For the recursive case,
			// the path is the result of filepath.Join, which calls filepath.Clean.
			// The initial case is not Cleaned, though, so we do this explicitly.
			//
			// This converts a path like "./io/" to "io". Without this step, running
			// "cd $GOROOT/src/pkg; go list ./io/..." would incorrectly skip the io
			// package, because prepending the prefix "./" to the unclean path would
			// result in "././io", and match("././io") returns false.
			path = filepath.Clean(path)
		}

		// Avoid directories we don't care about, but do not avoid "." or "..".
		_, elem := filepath.Split(path)
		dot := strings.HasPrefix(elem, ".") && elem != "." && elem != ".."
		if dot || strings.HasPrefix(elem, "_") {
			return filepath.SkipDir
		}

		name := prefix + filepath.ToSlash(path)
		if !match(name) {
			return nil
		}
		if _, err = build.ImportDir(path, 0); err != nil {
			if _, noGo := err.(*build.NoGoError); !noGo {
				log.Error("failed to read dir with err: %s", err)
			}
			return nil
		}
		pkgs = append(pkgs, name)
		return nil
	})

	if err != nil {
		log.Error("failed to walk the directories with err: %s", err)
	}

	return pkgs
}

// Find directory to begin the scan.
// Could be smarter but this one optimization
// is enough for now, since ... is usually at the
// end of a path.
func getStartDir(pattern string) string {
	var dir string

	i := strings.Index(pattern, "...")
	if i != -1 {
		dir, _ = path.Split(pattern[:i])
	} else {
		dir = pattern
	}

	return dir
}

func buildMatcher(pattern string) (string, func(name string) bool) {
	// pattern begins with ./ or ../.
	// path.Clean will discard the ./ but not the ../.
	// We need to preserve the ./ for pattern matching
	// and in the returned import paths.
	prefix := ""
	if strings.HasPrefix(pattern, "./") {
		prefix = "./"
	}
	match := matchPattern(pattern)

	return prefix, match
}

// matchPattern(pattern)(name) reports whether
// name matches pattern.  Pattern is a limited glob
// pattern in which '...' means 'any string' and there
// is no other special syntax.
func matchPattern(pattern string) func(name string) bool {
	re := regexp.QuoteMeta(pattern)
	re = strings.Replace(re, `\.\.\.`, `.*`, -1)

	// Special case: string ending in /
	if strings.HasSuffix(pattern, "/") {
		re = strings.TrimSuffix(re, `/`)
	}

	// Special case: foo/... matches foo too.
	if strings.HasSuffix(re, `/.*`) {
		re = re[:len(re)-len(`/.*`)] + `(/.*)?`
	}

	reg := regexp.MustCompile(`^` + re + `$`)
	return func(name string) bool {
		return reg.MatchString(name)
	}
}
