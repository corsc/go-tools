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

package main

import (
	"flag"
	"log"
	"os/exec"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	keyDirect   = "Direct"
	keyChild    = "Child"
	keyStdLib   = "Standard"
	keyVendored = "Vendored"
	keyExternal = "External"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalf("usage: tangled [directory name]")
	}

	directory := flag.Arg(0)
	deps := getDependencyList(directory)

	summary := buildSummary(deps)

	printSummary(summary)
}

func printSummary(in *stats) {
	fmt.Print("|---------------------------------------|\n")
	header := "| %-30s | %s |\n"
	fmt.Printf(header, "Count", "Type")
	fmt.Print("|---------------------------------------|\n")

	template := "| %-30s | %4d |\n"
	fmt.Printf(template, keyDirect, in.direct)
	fmt.Printf(template, keyChild, in.child)
	fmt.Printf(template, keyStdLib, in.stdLib)
	fmt.Printf(template, keyVendored, in.vendored)
	fmt.Printf(template, keyExternal, in.external)
	fmt.Print("|---------------------------------------|\n")
}

func buildSummary(deps *deps) *stats {
	out := &stats{}

	out.direct = len(deps.DirectImports)

	for _, thisDep := range deps.IndirectImports {
		if strings.HasPrefix(thisDep, "go/") {
			out.stdLib++
			continue
		}

		if strings.HasPrefix(thisDep, "vendor/golang_org/") {
			out.stdLib++
			continue
		}

		if !strings.Contains(thisDep, ".") {
			out.stdLib++
			continue
		}

		if strings.Contains(thisDep, "/vendor/") {
			out.vendored++
			continue
		}

		if strings.Contains(thisDep, deps.BasePath) {
			out.child++
			continue
		}

		out.external++
	}

	return out
}

func getDependencyList(directory string) *deps {
	bytes := goList(directory)

	out := &deps{}
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

type stats struct {
	direct   int
	child    int
	stdLib   int
	vendored int
	external int
}

// this is the JSON format returned by `go list --json`
type deps struct {
	BasePath string `json:"ImportPath"`
	DirectImports   []string `json:"Imports"`
	IndirectImports []string `json:"Deps"`
}
