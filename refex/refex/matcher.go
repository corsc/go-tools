package refex

import "strings"

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
		m.matches[index] = &match{
			startPos: regexMatch[0],
			endPos:   regexMatch[1],
			pattern:  regexp.String(),
		}
	}

	return nil
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
			loc := strings.Index(code, after)

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
	// TODO: add more regex conversion here
	chunk = strings.TrimPrefix(chunk, "(")
	chunk = strings.TrimSuffix(chunk, ")")

	for _, thisChar := range specialChars {
		chunk = strings.Replace(chunk, `\`+thisChar, thisChar, -1)
	}
	return chunk
}
