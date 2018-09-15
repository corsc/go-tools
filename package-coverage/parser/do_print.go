// Copyright (c) 2015-2018 Corey Scott (www.sage42.com), All Rights Reserved.
//
// NOTICE: All information contained herein is, and remains the property of Corey Scott.
// The intellectual and technical concepts contained herein are confidential, proprietary and controlled by Corey Scott
// and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this
// material in any way unless prior written permission is obtained from Corey Scott.
// Access to the source code contained herein is hereby forbidden to anyone except explicit written consent and subject
// to binding Confidentiality and Non-disclosure agreements explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of Corey Scott.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF COREY SCOTT IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS
// AND INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package parser

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/corsc/go-tools/package-coverage/config"
)

// DoPrint will output the coverage to StdOut
func DoPrint(cfg *config.Config, path string, exclusions *regexp.Regexp) bool {
	coverageOk := true

	if !cfg.DoPrint {
		return coverageOk
	}

	buffer := bytes.Buffer{}
	if cfg.SingleDir {
		coverageOk = PrintCoverageSingle(&buffer, path, cfg.MinCoverage, cfg.Prefix, cfg.Depth)
	} else {
		coverageOk = PrintCoverage(&buffer, path, exclusions, cfg.MinCoverage, cfg.Prefix, cfg.Depth)
	}

	fmt.Print(buffer.String())
	return coverageOk
}
