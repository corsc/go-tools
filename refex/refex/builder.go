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
