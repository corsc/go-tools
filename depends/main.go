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

	"github.com/corsc/go-tools/depends/internal"
)

func main() {
	cfg := &config{}
	setUsage(cfg)
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}

	directory := flag.Arg(0)
	focusPackage := flag.Arg(1)
	masterList := internal.GetDependantsList(directory)

	summary := internal.BuildSummary(focusPackage, masterList)

	printOutput(cfg, summary)
}

func printOutput(cfg *config, in *internal.Summary) {
	if cfg.listCSV {
		internal.PrintCSVList(in)
		return
	}

	if cfg.listDirect {
		internal.PrintDirect(in)
		return
	}
	if cfg.listTest {
		internal.PrintTest(in)
		return
	}

	internal.PrintFullList(in)
}

type config struct {
	listDirect bool
	listTest   bool
	listCSV    bool
}

func setUsage(cfg *config) {
	flag.BoolVar(&cfg.listDirect, "direct", false, "list only direct dependencies")
	flag.BoolVar(&cfg.listTest, "test", false, "list only test dependencies")
	flag.BoolVar(&cfg.listCSV, "csv", false, "print list as CSV (overrides other options)")

	flag.Usage = func() {
		fmt.Print("Usage of depends:\n")
		fmt.Printf("\tdepends [flags] [base directory] [pkg]\n")
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}
}
