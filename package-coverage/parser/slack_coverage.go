package parser

import (
	"fmt"
	"log"
	"regexp"

	"bytes"
	"net/http"

	"strings"

	"io/ioutil"

	"github.com/corsc/go-tools/package-coverage/utils"
)

// SlackCoverageSingle is the same as GetCoverage only for 1 directory only
func SlackCoverageSingle(path string, matcher *regexp.Regexp, webhook string, prefix string, depth int) {
	var fullPath string
	if path == "./" {
		fullPath = utils.GetCurrentDir()
	} else {
		fullPath = utils.GetCurrentDir() + path + "/"
	}
	fullPath += "profile.cov"

	pkgs, coverageData := getCoverageData([]string{fullPath}, matcher)
	prepareAndSendToSlack(pkgs, coverageData, webhook, prefix, depth)
}

// SlackCoverage will attempt to calculate and print the coverage from the supplied coverage files
func SlackCoverage(basePath string, matcher *regexp.Regexp, webhook string, prefix string, depth int) {
	paths, err := utils.FindAllCoverageFiles(basePath)
	if err != nil {
		log.Panicf("error file finding coverage files %s", err)
	}

	pkgs, coverageData := getCoverageData(paths, matcher)
	prepareAndSendToSlack(pkgs, coverageData, webhook, prefix, depth)
}

func prepareAndSendToSlack(pkgs []string, coverageData coverageByPackage, webhook string, prefix string, depth int) {
	lines := 0
	output := ""

	for _, pkg := range pkgs {
		cover := coverageData[pkg]
		covered, stmts := getStats(cover)

		pkgFormatted := strings.Replace(pkg, prefix, "", -1)
		pkgDepth := strings.Count(pkgFormatted, "/")

		if depth > 0 && pkgDepth <= depth {
			var color string
			if covered > 70 {
				color = "good"
			} else if covered > 50 {
				color = "warning"
			} else {
				color = "danger"
			}

			if lines > 0 {
				output += ","
			}

			output += fmt.Sprintf("{ \"color\": \"%s\", \"text\": \"%-50s %3.2f%% (%0.0f statements)\" }", color, pkgFormatted, covered, stmts)
			if lines >= 20 {
				sendToSlack(webhook, output)
				lines = 0
				output = ""
			} else {
				lines++
			}
		}
	}

	if len(output) > 0 {
		sendToSlack(webhook, output)
	}
}

func sendToSlack(webhook string, attachments string) {
	message := `{ "username": "Test Coverage Bot", "attachments": [ ` + attachments + ` ] }`

	req, err := http.NewRequest("POST", webhook, bytes.NewBufferString(message))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		panic(fmt.Sprintf("unexpected response code %d; body: %s; payload: %s", resp.StatusCode, body, message))
	}
}
