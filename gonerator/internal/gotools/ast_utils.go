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

package gotools

import (
	"go/token"
)

const lineBreak = '\n'

// GetLineBoundary will return the start and end position of a line containing position pos
// NOTE: this method with panic if an invalid position is supplied
func GetLineBoundary(source []byte, pos token.Pos) (int, int) {
	// Run to end of line in both directions if not at line start/end.
	charAtPos := source[pos]
	if charAtPos == lineBreak {
		return int(pos), int(pos)
	}

	startPos, endPos := int(pos), int(pos)+1
	for startPos > 0 && source[startPos-1] != lineBreak {
		startPos--
	}

	for endPos < len(source) && source[endPos-1] != lineBreak {
		endPos++
	}

	return startPos, endPos
}
