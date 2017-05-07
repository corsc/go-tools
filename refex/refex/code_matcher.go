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
	"strings"
)

type codeMatcher interface {
	match(code string, pattern pattern) ([]*match, error)
}

type codeMatcherImpl struct {
	code    string
	matches []*match
}

func (m *codeMatcherImpl) match(code string, pattern pattern) ([]*match, error) {
	err := m.find(code, pattern)
	if err != nil {
		return nil, err
	}

	err = m.buildParts()
	if err != nil {
		return nil, err
	}

	return m.matches, nil
}

func (m *codeMatcherImpl) find(code string, pattern pattern) error {
	m.code = code

	regexp, err := pattern.regexp()
	if err != nil {
		return err
	}

	regexMatches := regexp.FindAllStringIndex(code, -1)
	if regexMatches == nil {
		// not an error but still nothing more to do
		return nil
	}

	m.matches = make([]*match, len(regexMatches))
	for index, regexMatch := range regexMatches {
		startPos := regexMatch[0]
		endPos := regexMatch[1]

		endPos -= m.adjustForDoubleBrackets(code, startPos, endPos)
		endPos -= m.adjustForGreedy(code, startPos, endPos)

		m.matches[index] = &match{
			startPos: startPos,
			endPos:   endPos,
			pattern:  regexp.String(),
		}
	}

	return nil
}

// Adjust for cases where the params contain `()` or `(something)`
func (m *codeMatcherImpl) adjustForDoubleBrackets(code string, startPos int, endPos int) int {
	adjustment := 0

	codeMatch := code[startPos:endPos]
	if leftBrackets := strings.Count(codeMatch, "("); leftBrackets > 0 {
		lastRightBracket := strings.LastIndex(codeMatch, ")")
		if lastRightBracket > 0 {
			adjustment = (len(codeMatch) - lastRightBracket - 1)
		}
	}

	return adjustment
}

// Adjust for cases where Greedy matching matches too many closing brackets
func (m *codeMatcherImpl) adjustForGreedy(code string, startPos int, endPos int) int {
	adjustment := 0

	codeMatch := code[startPos:endPos]
	leftBrackets := strings.Count(codeMatch, "(")
	rightBrackets := strings.Count(codeMatch, ")")

	if rightBrackets > leftBrackets {
		lastRightBracket := strings.LastIndex(codeMatch, ")")
		if lastRightBracket > 0 {
			adjustment = (len(codeMatch) - lastRightBracket)
		}
	}

	return adjustment
}

func (m *codeMatcherImpl) buildParts() error {
	if len(m.matches) == 0 {
		// not an error but still nothing more to do
		return nil
	}

	lastPos := 0
	for _, match := range m.matches {
		match.parts = []*part{}

		before := m.code[lastPos:match.startPos]
		code := m.code[match.startPos:match.endPos]

		patternChunks := strings.Split(match.pattern, wildcard)

		for index, chunkRaw := range patternChunks {
			chunk := m.extractCodeFromPattern(chunkRaw)
			before += chunk
			if index == 0 {
				match.parts = append(match.parts, &part{code: before})
			}

			code = strings.TrimPrefix(code, chunk)

			nextIndex := index + 1
			if nextIndex > (len(patternChunks) - 1) {
				continue
			}

			after := m.extractCodeFromPattern(patternChunks[nextIndex])

			// grab the first occurrence or the last depending where we are in the pattern
			var loc int
			if nextIndex == (len(patternChunks) - 1) {
				loc = strings.LastIndex(code, after)
			} else {
				loc = strings.Index(code, after)
			}

			thisFragment := code[:loc]

			match.parts = append(match.parts, &part{
				code:  thisFragment,
				isArg: true,
				index: (index + 1),
			})
			match.parts = append(match.parts, &part{code: after})

			code = code[loc:]
		}

		lastPos = match.endPos
	}

	return nil
}

func (m *codeMatcherImpl) extractCodeFromPattern(chunk string) string {
	chunk = strings.TrimPrefix(chunk, patternPrefix)
	chunk = strings.TrimSuffix(chunk, patternSuffix)

	for _, thisChar := range specialChars {
		chunk = strings.Replace(chunk, `\`+thisChar, thisChar, -1)
	}
	return chunk
}
