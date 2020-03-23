// Copyright (c) 2015-2018 Corey Scott (www.sage42.com), All Rights Reserved.
//
// NOTICE: All information contained herein is, and remains the property of Corey Scott.
// The intellectual and technical concepts contained herein are confidential, proprietary and controlled by Corey Scott
// and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this
// material in any way unless prior written permission is obtained from Corey Scott.
// Access to the source code contained herein is hereby forbidden to anyone except explicit written consent and subject
// to binding Confidentiality and Non-disclosure agreements explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of Corey Scott.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF COREY SCOTT IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS
// AND INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package config

import (
	"flag"
	"os"
	"runtime"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// Config for this application
type Config struct {
	// Verbose mode is useful for debugging this tool (default:false)
	Verbose bool

	// Concurrency used to run tests (default: runtime.NumCPU())
	Concurrency int

	// Quiet mode will suppress the StdOut messages from go test (default:true)
	Quiet bool

	// DoAll is short form/convenience method for -c -p -d (calculate, output and clean up)
	DoAll bool

	// Coverage controls is coverage is calculator (or reused from previous run)
	Coverage bool

	// SingleDir will only generate for the supplied directory (no recursion and will ignore -i)
	SingleDir bool

	// DoClean will "clean up" by removing any calculated coverage files
	DoClean bool

	// DoPrint will output the result to StdOut
	DoPrint bool

	// IgnorePaths allows you to ignore file paths matching the specified regex
	// (match directories by surrounding the directory name with slashes; match files by prefixing with a slash)
	IgnorePaths string

	// WebHook is the Slack WebHook URL (missing means don't send)
	WebHook string

	// ChannelOverride allows you to override the WebHook's default Slack channel
	ChannelOverride string

	// Prefix is the directory structure to be removed from all package names (makes the output cleaner)
	Prefix string

	// Depth is how many levels of coverage to output (default is 0 = all)
	Depth int

	// MinCoverage causes output to StdOut to be colored red for any package below this amount of coverage
	MinCoverage int

	// Tags is the go build tags to be added in go test calls
	Tags string

	// Race is used to enable --race flag
	Race bool
}

// GetConfig will extra config from flags and return
func GetConfig() *Config {
	cfg := &Config{}

	// extract from flags
	flag.BoolVar(&(cfg.Verbose), "v", false, "verbose mode is useful for debugging this tool")
	flag.IntVar(&(cfg.Concurrency), "concurrency", runtime.NumCPU(), "concurrency used to run tests (default: runtime.NumCPU())")
	flag.BoolVar(&(cfg.Quiet), "q", true, "quiet mode will suppress the stdOut messages from go test")
	flag.BoolVar(&(cfg.Coverage), "c", false, "generate coverage")
	flag.BoolVar(&(cfg.SingleDir), "s", false, "only generate for the supplied directory (no recursion / will ignore -i)")
	flag.BoolVar(&(cfg.DoClean), "d", false, "clean")
	flag.BoolVar(&(cfg.DoPrint), "p", false, "print coverage to stdout")
	flag.StringVar(&(cfg.IgnorePaths), "i", `./\.git.*|./_.*`, "ignore file paths matching the specified regex (match directories by surrounding the directory name with slashes; match files by prefixing with a slash)")
	flag.StringVar(&(cfg.WebHook), "webhook", "", "Slack webhook URL (missing means don't send)")
	flag.StringVar(&(cfg.ChannelOverride), "channel", "", "Slack channel (missing means use the default channel for this webhook)")
	flag.StringVar(&(cfg.Prefix), "prefix", "", "prefix is the directory structure to be removed from all package names (makes the output cleaner)")
	flag.IntVar(&(cfg.Depth), "depth", 0, "How many levels of coverage to output (default is 0 = all)")
	flag.IntVar(&(cfg.MinCoverage), "m", 0, "minimum coverage")
	flag.StringVar(&(cfg.Tags), "tags", ``, "go build tags to be added in go test calls")
	flag.BoolVar(&(cfg.Race), "r", false, "enable race detection during testing")
	flag.BoolVar(&(cfg.DoAll), "a", true, "short form/convenience method for -c -p -d (calculate, output and clean up)")
	flag.Parse()

	// initialize verbose mode
	if cfg.Verbose {
		utils.LogWhenVerbose("Config: %#v", cfg)
	} else {
		utils.VerboseOff()
	}

	// validate config
	if cfg.Depth > 0 && len(cfg.Prefix) == 0 {
		println("You must specify a prefix when using -depth")
		os.Exit(-1)
	}

	// Set "default" mode (Calculate+Print+Clean up) when selected
	if cfg.DoAll {
		cfg.Coverage = true
		cfg.DoPrint = true
		cfg.DoClean = true
	}

	return cfg
}
