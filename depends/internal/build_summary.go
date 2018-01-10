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

// BuildSummary will summarize dependencies by type
func BuildSummary(focusPackage string, masterList *MasterList) *Summary {
	out := &Summary{
		direct: []string{},
		test:   []string{},
	}

	for _, thisPkg := range masterList.Pkgs {
		if thisPkg.BasePath == focusPackage {
			// skip the package we are looking at
			continue
		}

		for _, thisImport := range thisPkg.DirectImports {
			if thisImport == focusPackage {
				out.direct = append(out.direct, thisPkg.BasePath)
			}
		}

		for _, thisImport := range thisPkg.TestImports {
			if thisImport == focusPackage {
				out.direct = append(out.test, thisPkg.BasePath)
			}
		}
	}

	return out
}

// Summary is a dependency summary
type Summary struct {
	direct []string
	test   []string
}
