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

package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/corsc/go-commons/iocloser"
	"github.com/corsc/go-tools/package-coverage/utils"
)

// SlackCoverage will attempt to calculate and output the coverage from the supplied coverage files to Slack
func SlackCoverage(basePath string, exclusionsMatcher *regexp.Regexp, webHook string, channelOverride string, prefix string, depth int) {
	paths, err := utils.FindAllCoverageFiles(basePath)
	if err != nil {
		log.Panicf("error file finding coverage files %s", err)
	}

	pkgs, coverageData := getCoverageData(paths, exclusionsMatcher)
	prepareAndSendToSlack(pkgs, coverageData, webHook, channelOverride, prefix, depth)
}

// SlackCoverageSingle is the same as SlackCoverage only for 1 directory only
func SlackCoverageSingle(path string, webHook string, channelOverride string, prefix string, depth int) {
	var fullPath string
	if path == "./" {
		fullPath = utils.GetCurrentDir()
	} else {
		fullPath = utils.GetCurrentDir() + path + "/"
	}
	fullPath += "profile.cov"

	contents := getFileContents(fullPath)
	pkgs, coverageData := getCoverageByContents(contents)

	prepareAndSendToSlack(pkgs, coverageData, webHook, channelOverride, prefix, depth)
}

// prepare the slack message format and send.
// Notes:
// * the message uses the Slack message "attachments"; one attachment per package
// * each package is prefixed with a color highlight that corresponds to coverage amounts.
// (coverage > 70% is green.  70% > x >= 50 is orange.  coverage < 50% is red)
func prepareAndSendToSlack(pkgs []string, coverageData coverageByPackage, webhook string, channelOverride string, prefix string, depth int) {
	lines := 0
	output := ""

	for _, pkg := range pkgs {
		cover := coverageData[pkg]
		covered, _, statements := getSummaryValues(cover)

		pkgFormatted := strings.Replace(pkg, prefix, "", -1)
		pkgDepth := strings.Count(pkgFormatted, "/")

		if depth > 0 {
			if pkgDepth <= depth {
				addLineSlack(&output, pkgFormatted, covered, statements, lines)
				if lines >= 18 {
					sendToSlack(webhook, channelOverride, output)
					lines = 0
					output = ""
				} else {
					lines++
				}
			}
		} else {
			addLineSlack(&output, pkgFormatted, covered, statements, 0)
		}
	}

	if len(output) > 0 {
		sendToSlack(webhook, channelOverride, output)
	}
}

func addLineSlack(output *string, pkgFormatted string, covered float64, stmts float64, lines int) {
	var color string
	if covered > 70 {
		color = "good"
	} else if covered > 50 {
		color = "warning"
	} else {
		color = "danger"
	}

	if lines > 0 {
		*output += ","
	}

	*output += fmt.Sprintf("{ \"color\": \"%s\", \"text\": \"%-50s %3.2f%% (%0.0f statements)\" }", color, pkgFormatted, covered, stmts)
}

// call the Slack incoming webHook API to send the message
func sendToSlack(webHook string, channelOverride string, attachments string) {
	customChannel := ""
	if len(channelOverride) > 0 {
		customChannel = `, "channel": "` + channelOverride + `"`
	}

	message := `{ "username": "Test Coverage Bot", "attachments": [ ` + attachments + ` ] ` + customChannel + ` }`

	resp, err := http.Post(webHook, "application/json", bytes.NewBufferString(message))
	if err != nil {
		panic(err)
	}
	defer iocloser.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Panicf("unexpected response code %d; body: %s; payload: %s", resp.StatusCode, body, message)
	}
}
