package refex

import (
	"errors"
)

type codeBuilder interface {
	build(beforeParts []*part, afterParts []*part) (string, error)
}

type codeBuilderImpl struct{}

func (c *codeBuilderImpl) build(beforeParts []*part, afterParts []*part) (string, error) {
	if len(afterParts) > len(beforeParts) {
		return "", errors.New("number of after parts cannot be more than before parts")
	}

	newCode := ""
	for _, afterPart := range afterParts {
		if afterPart.isArg {
			newCode += c.buildArg(afterPart, beforeParts)
		} else {
			newCode += c.buildCode(afterPart)
		}
	}

	return newCode, nil
}

func (c *codeBuilderImpl) buildArg(afterPart *part, beforeParts []*part) string {
	for _, beforePart := range beforeParts {
		if beforePart.isArg && beforePart.index == afterPart.index {
			return beforePart.code
		}
	}
	return ""
}

func (c *codeBuilderImpl) buildCode(afterPart *part) string {
	return afterPart.code
}
