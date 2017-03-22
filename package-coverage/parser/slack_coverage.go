package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// SlackCoverage will attempt to calculate and output the coverage from the supplied coverage files to Slack
func SlackCoverage(basePath string, exclusionsMatcher *regexp.Regexp, webHook string, prefix string, depth int) {
	paths, err := utils.FindAllCoverageFiles(basePath)
	if err != nil {
		log.Panicf("error file finding coverage files %s", err)
	}

	pkgs, coverageData := getCoverageData(paths, exclusionsMatcher)
	prepareAndSendToSlack(pkgs, coverageData, webHook, prefix, depth)
}

// SlackCoverageSingle is the same as SlackCoverage only for 1 directory only
func SlackCoverageSingle(path string, webHook string, prefix string, depth int) {
	var fullPath string
	if path == "./" {
		fullPath = utils.GetCurrentDir()
	} else {
		fullPath = utils.GetCurrentDir() + path + "/"
	}
	fullPath += "profile.cov"

	contents := getFileContents(fullPath)
	pkgs, coverageData := getCoverageByContents(contents)

	prepareAndSendToSlack(pkgs, coverageData, webHook, prefix, depth)
}

// prepare the slack message format and send.
// Notes:
// * the message uses the Slack message "attachments"; one attachment per package
// * each package is prefixed with a color highlight that corresponds to coverage amounts.
// (coverage > 70% is green.  70% > x >= 50 is orange.  coverage < 50% is red)
func prepareAndSendToSlack(pkgs []string, coverageData coverageByPackage, webhook string, prefix string, depth int) {
	lines := 0
	output := ""

	for _, pkg := range pkgs {
		cover := coverageData[pkg]
		covered, statements := getStats(cover)

		pkgFormatted := strings.Replace(pkg, prefix, "", -1)
		pkgDepth := strings.Count(pkgFormatted, "/")

		if depth > 0 {
			if pkgDepth <= depth {
				addLineSlack(&output, pkgFormatted, covered, statements, lines)
				if lines >= 18 {
					sendToSlack(webhook, output)
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
		sendToSlack(webhook, output)
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
func sendToSlack(webHook string, attachments string) {
	message := `{ "username": "Test Coverage Bot", "attachments": [ ` + attachments + ` ] }`

	resp, err := http.Post(webHook, "application/json", bytes.NewBufferString(message))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Panicf("unexpected response code %d; body: %s; payload: %s", resp.StatusCode, body, message)
	}
}
