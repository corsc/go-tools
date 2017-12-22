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
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
)

// GetDependencyList returns the dependency list of given package (directory)
func GetDependencyList(directory string) *Deps {
	bytes := goList(directory)

	out := &Deps{}
	err := json.Unmarshal(bytes, out)
	if err != nil {
		log.Fatalf("failed to parse go list data with err %s", err)
	}

	return out
}

func goList(directory string) []byte {
	cmd := exec.Command("go", "list", "--json")
	cmd.Dir = directory

	output := &bytes.Buffer{}
	catchErr := &bytes.Buffer{}

	cmd.Stdout = output
	cmd.Stderr = catchErr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to get deps from go list with err %s", err)
	}

	if catchErr.Len() > 0 {
		log.Fatalf("failed to get deps from go list with err %s", err)
	}

	return output.Bytes()
}

// Deps is the JSON format returned by `go list --json`
type Deps struct {
	BasePath        string   `json:"ImportPath"`
	DirectImports   []string `json:"Imports"`
	IndirectImports []string `json:"Deps"`
}
