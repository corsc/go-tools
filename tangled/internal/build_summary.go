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

package internal

import (
	"strings"
)

// BuildSummary will summarize dependencies by type
func BuildSummary(deps *Deps) *Summary {
	out := &Summary{
		Pkg:      deps.BasePath,
		direct:   map[string]struct{}{},
		child:    map[string]struct{}{},
		stdLib:   map[string]struct{}{},
		vendored: map[string]struct{}{},
		external: map[string]struct{}{},
	}

	for _, thisDep := range deps.DirectImports {
		out.direct[thisDep] = struct{}{}
	}

	for _, thisDep := range deps.IndirectImports {
		if strings.HasPrefix(thisDep, "go/") {
			out.stdLib[thisDep] = struct{}{}
			continue
		}

		if strings.HasPrefix(thisDep, "vendor/golang_org/") {
			out.stdLib[thisDep] = struct{}{}
			continue
		}

		if !strings.Contains(thisDep, ".") {
			out.stdLib[thisDep] = struct{}{}
			continue
		}

		if strings.Contains(thisDep, "/vendor/") {
			out.vendored[thisDep] = struct{}{}
			continue
		}

		if strings.Contains(thisDep, deps.BasePath) {
			out.child[thisDep] = struct{}{}
			continue
		}

		out.external[thisDep] = struct{}{}
	}

	return out
}

// Summary is a dependency summary
type Summary struct {
	Pkg      string
	direct   map[string]struct{}
	child    map[string]struct{}
	stdLib   map[string]struct{}
	vendored map[string]struct{}
	external map[string]struct{}
}
