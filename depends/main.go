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
	masterList := internal.GetDependantsList(directory)

	summary := internal.BuildSummary(masterList)

	printOutput(cfg, summary)
}

func printOutput(cfg *config, in *internal.Summary) {
	if cfg.list {
		internal.PrintFullList(in)
		return
	}

	if cfg.summaryCSV {
		internal.PrintSummaryCSV(in)
		return
	}

	internal.PrintSummary(in)

	if cfg.listChild {
		internal.PrintChild(in)
	}
	if cfg.listStdLib {
		internal.PrintStdLib(in)
	}
	if cfg.listExternal {
		internal.PrintExternal(in)
	}
	if cfg.listVendored {
		internal.PrintVendored(in)
	}
}

type config struct {
	list         bool
	listChild    bool
	listStdLib   bool
	listVendored bool
	listExternal bool
	summaryCSV   bool
}

func setUsage(cfg *config) {
	flag.BoolVar(&cfg.list, "list", false, "list everything ")
	flag.BoolVar(&cfg.listChild, "child", false, "list child dependencies")
	flag.BoolVar(&cfg.listStdLib, "std", false, "list standard library dependencies")
	flag.BoolVar(&cfg.listVendored, "vendored", false, "list vendored dependencies")
	flag.BoolVar(&cfg.listExternal, "external", false, "list external dependencies")
	flag.BoolVar(&cfg.summaryCSV, "scsv", false, "print summary as CSV (overrides other options)")
}
