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
	"fmt"
	"os"
	"regexp"

	"github.com/corsc/go-tools/package-coverage/config"
	"github.com/corsc/go-tools/package-coverage/generator"
	"github.com/corsc/go-tools/package-coverage/parser"
	"github.com/corsc/go-tools/package-coverage/utils"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %s\n", r)
		}
	}()

	// get config and environment
	cfg := config.GetConfig()
	startDir := utils.GetCurrentDir()
	path := getPath()

	// build exclusions regex
	var exclusions *regexp.Regexp
	if cfg.IgnorePaths != "" {
		exclusions = regexp.MustCompile(cfg.IgnorePaths)
	}

	// calculate coverage
	generator.Calculate(cfg, path, exclusions)

	// switch back to start dir
	err := os.Chdir(startDir)
	if err != nil {
		panic(err)
	}

	// output coverage to StdOut
	coverageOk := parser.DoPrint(cfg, path, exclusions)

	// output to Slack
	parser.DoSlack(cfg, path, exclusions)

	// clean up
	generator.DoClean(cfg, path, exclusions)

	// signal success or not
	if !coverageOk {
		os.Exit(-1)
	}
}

func getPath() string {
	path := flag.Arg(0)
	if path == "" {
		println("Please include a directory as the last argument")
		os.Exit(-1)
	}
	return path
}
