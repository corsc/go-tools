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
	"sort"
)

const (
	keyDirect   = "Direct"
	keyChild    = "Child"
	keyStdLib   = "Std Lib"
	keyExternal = "External"
	keyVendored = "Vendored"
)

func main() {
	cfg := &config{}
	setUsage(cfg)
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}

	directory := flag.Arg(0)
	deps := getDependencyList(directory)

	summary := buildSummary(deps)

	printOutput(cfg, summary)
}

func printOutput(cfg *config, in *stats) {
	printSummary(in)

	if cfg.listDirect {
		printList(keyDirect, in.direct)
	}
	if cfg.listChild {
		printList(keyChild, in.child)
	}
	if cfg.listStdLib {
		printList(keyStdLib, in.stdLib)
	}
	if cfg.listExternal {
		printList(keyExternal, in.external)
	}
	if cfg.listVendored {
		printList(keyVendored, in.vendored)
	}
}

func printList(title string, items map[string]struct{}) {
	sortedItems := make([]string, 0, len(items))
	for key := range items {
		sortedItems = append(sortedItems, key)
	}
	sort.Strings(sortedItems)

	header := "\n%-30s\n"
	fmt.Printf(header, title)
	fmt.Print("------------------------------\n")

	template := "%s\n"
	for _, item := range sortedItems {
		fmt.Printf(template, item)
	}
	println()
}

func printSummary(in *stats) {
	fmt.Print("|---------------------------------------|\n")
	header := "| %-30s | %s |\n"
	fmt.Printf(header, "Count", "Type")
	fmt.Print("|---------------------------------------|\n")

	template := "| %-30s | %4d |\n"
	fmt.Printf(template, keyDirect, len(in.direct))
	fmt.Printf(template, keyChild, len(in.child))
	fmt.Printf(template, keyStdLib, len(in.stdLib))
	fmt.Printf(template, keyExternal, len(in.external))
	fmt.Printf(template, keyVendored, len(in.vendored))
	fmt.Print("|---------------------------------------|\n")
}

func buildSummary(deps *deps) *stats {
	out := &stats{
		direct  : map[string]struct{}{},
		child   : map[string]struct{}{},
		stdLib  : map[string]struct{}{},
		vendored : map[string]struct{}{},
		external : map[string]struct{}{},
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
	direct   map[string]struct{}
	child    map[string]struct{}
	stdLib   map[string]struct{}
	vendored map[string]struct{}
	external map[string]struct{}
}

// this is the JSON format returned by `go list --json`
type deps struct {
	BasePath        string   `json:"ImportPath"`
	DirectImports   []string `json:"Imports"`
	IndirectImports []string `json:"Deps"`
}

type config struct {
	listDirect   bool
	listChild    bool
	listStdLib   bool
	listVendored bool
	listExternal bool
}

func setUsage(cfg *config) {
	flag.BoolVar(&cfg.listDirect, "direct", false, "list direct dependencies")
	flag.BoolVar(&cfg.listChild, "child", false, "list child dependencies")
	flag.BoolVar(&cfg.listStdLib, "std", false, "list standard library dependencies")
	flag.BoolVar(&cfg.listVendored, "vendored", false, "list vendored dependencies")
	flag.BoolVar(&cfg.listExternal, "external", false, "list external dependencies")
}
