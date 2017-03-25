package refex

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type pattern interface {
	build(in string) ([]*part, error)
	regexp() (*regexp.Regexp, error)
}

type patternImpl struct {
	pattern string
	parts   []*part
}

func (p *patternImpl) build(in string) ([]*part, error) {
	p.pattern = in

	matches := regexp.MustCompile(`(\$)[0-9]+(\$)`).FindAllStringIndex(in, -1)
	if matches == nil {
		return nil, fmt.Errorf("no tokens found in supplied transform '%s'", in)
	}

	lastPos := 0
	for _, match := range matches {
		start := in[lastPos:match[0]]
		middle := in[(match[0] + 1):(match[1] - 1)]

		p.parts = append(p.parts, &part{code: start})

		index, _ := strconv.Atoi(middle)
		p.parts = append(p.parts, &part{isArg: true, index: index})

		lastPos = match[1]
	}

	end := in[lastPos:]
	p.parts = append(p.parts, &part{code: end})

	if len(p.parts) == 0 {
		return nil, fmt.Errorf("failed to find any parts in transform '%s'", in)
	}

	return p.parts, nil
}

func (p *patternImpl) regexp() (*regexp.Regexp, error) {
	regex := strings.Replace(p.pattern, `(`, `\(`, -1)
	regex = strings.Replace(regex, `)`, `\)`, -1)

	for _, part := range p.parts {
		regex = strings.Replace(regex, "$"+strconv.Itoa(part.index)+"$", ")(.*)(", 1)
	}

	return regexp.Compile(`(` + regex + `)`)
}
