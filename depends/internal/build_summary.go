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
func BuildSummary(masterList *MasterList) *Summary {
	out := &Summary{
		child:    map[string]*SummaryItem{},
		stdLib:   map[string]*SummaryItem{},
		vendored: map[string]*SummaryItem{},
		external: map[string]*SummaryItem{},
	}

	for _, thisPkg := range masterList.Pkgs {
		basePath := thisPkg.BasePath

		for _, thisDep := range thisPkg.DirectImports {
			sortItem(basePath, thisDep, out)
		}

		for _, thisDep := range thisPkg.IndirectImports {
			sortItem(basePath, thisDep, out)
		}
	}

	return out
}

func sortItem(basePath string, thisDep string, out *Summary) {
	if strings.HasPrefix(thisDep, "go/") {
		out.getChildItem(out.stdLib, thisDep).addDependent(basePath)
		return
	}

	if strings.HasPrefix(thisDep, "vendor/golang_org/") {
		out.getChildItem(out.stdLib, thisDep).addDependent(basePath)
		return
	}

	if !strings.Contains(thisDep, ".") {
		out.getChildItem(out.stdLib, thisDep).addDependent(basePath)
		return
	}

	if strings.Contains(thisDep, "/vendor/") {
		out.getChildItem(out.vendored, thisDep).addDependent(basePath)
		return
	}

	if strings.Contains(thisDep, basePath) {
		out.getChildItem(out.child, thisDep).addDependent(basePath)
		return
	}

	out.getChildItem(out.external, thisDep).addDependent(basePath)
}

// Summary is a dependency summary
type Summary struct {
	child    map[string]*SummaryItem
	external map[string]*SummaryItem
	stdLib   map[string]*SummaryItem
	vendored map[string]*SummaryItem
}

func (s *Summary) getChildItem(m map[string]*SummaryItem, pkg string) *SummaryItem {
	out, found := m[pkg]
	if !found {
		m[pkg] = &SummaryItem{
			Pkg:        pkg,
			Dependents: map[string]struct{}{},
		}
		out = m[pkg]
	}
	return out
}

// SummaryItem is a dependency summary
type SummaryItem struct {
	Pkg        string
	Dependents map[string]struct{}
}

func (s *SummaryItem) addDependent(dep string) {
	s.Dependents[dep] = struct{}{}
}
