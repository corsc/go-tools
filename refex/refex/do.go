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

package refex

import (
	"io/ioutil"
)

// DoFile wraps Do() and loads the code from the filename supplied
func DoFile(fileName string, before string, after string) (string, error) {
	codeIn, err := getFileContents(fileName)
	if err != nil {
		return "", err
	}

	return Do(codeIn, before, after)
}

// Do will replace all matches of "before" with "after" in "codeIn"
func Do(codeIn string, before string, after string) (string, error) {
	output := ""

	// build a pattern to be replaced
	var beforePattern pattern = &patternImpl{}
	_, err := beforePattern.build(before)
	if err != nil {
		return "", err
	}

	// find that pattern in the code
	var codeMatcher codeMatcher = &codeMatcherImpl{}
	beforeMatches, err := codeMatcher.match(codeIn, beforePattern)
	if err != nil {
		return "", err
	}

	// build a pattern for the replacement
	var afterPattern pattern = &patternImpl{}
	afterParts, err := afterPattern.build(after)
	if err != nil {
		return "", err
	}

	// build new code from matches and after pattern
	var codeBuilder codeBuilder = &codeBuilderImpl{}

	lastPos := 0
	for _, beforeMatch := range beforeMatches {
		newCode, err := codeBuilder.build(beforeMatch.parts, afterParts)
		if err != nil {
			return "", err
		}

		// replace old code with new
		output += codeIn[lastPos:beforeMatch.startPos]
		output += newCode

		lastPos = beforeMatch.endPos
	}

	if len(codeIn) > lastPos {
		output += codeIn[lastPos:]
	}

	return output, nil
}

func getFileContents(fileName string) (string, error) {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}
