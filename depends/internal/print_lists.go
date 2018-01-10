// Copyright (c) 2012-2017 Grab Taxi Holdings PTE LTD (GRAB), All Rights Reserved. NOTICE: All information contained herein
// is, and remains the property of GRAB. The intellectual and technical concepts contained herein are confidential, proprietary
// and controlled by GRAB and may be covered by patents, patents in process, and are protected by trade secret or copyright law.
//
// You are strictly forbidden to copy, download, store (in any medium), transmit, disseminate, adapt or change this material
// in any way unless prior written permission is obtained from GRAB. Access to the source code contained herein is hereby
// forbidden to anyone except current GRAB employees or contractors with binding Confidentiality and Non-disclosure agreements
// explicitly covering such access.
//
// The copyright notice above does not evidence any actual or intended publication or disclosure of this source code,
// which includes information that is confidential and/or proprietary, and is a trade secret, of GRAB.
//
// ANY REPRODUCTION, MODIFICATION, DISTRIBUTION, PUBLIC PERFORMANCE, OR PUBLIC DISPLAY OF OR THROUGH USE OF THIS SOURCE
// CODE WITHOUT THE EXPRESS WRITTEN CONSENT OF GRAB IS STRICTLY PROHIBITED, AND IN VIOLATION OF APPLICABLE LAWS AND
// INTERNATIONAL TREATIES. THE RECEIPT OR POSSESSION OF THIS SOURCE CODE AND/OR RELATED INFORMATION DOES NOT CONVEY
// OR IMPLY ANY RIGHTS TO REPRODUCE, DISCLOSE OR DISTRIBUTE ITS CONTENTS, OR TO MANUFACTURE, USE, OR SELL ANYTHING
// THAT IT MAY DESCRIBE, IN WHOLE OR IN PART.

package internal

import (
	"fmt"
)

const (
	keyDirect = "Direct"
	keyTest   = "Test"
)

// PrintFullList of all packages and their dependents
func PrintFullList(in *Summary) {
	PrintDirect(in)
	PrintTest(in)
}

func doPrintFullList(title string, itemsMap []string) {
	header := "----------------------------\n%-30s\n"
	fmt.Printf(header, title)

	for _, item := range itemsMap {
		template := "    %-60s\n"
		fmt.Printf(template, item)
	}
	println()
}

// PrintDirect will output the direct deps
func PrintDirect(in *Summary) {
	doPrintFullList(keyDirect, in.direct)
}

// PrintTest will output the std lib deps
func PrintTest(in *Summary) {
	doPrintFullList(keyTest, in.test)
}

// PrintFullList of all packages and their dependents
func PrintCSVList(in *Summary) {
	doPrintCSVList(keyDirect, in.direct)
	doPrintCSVList(keyTest, in.test)
}

func doPrintCSVList(title string, items []string) {
	for _, item := range items {
		template := "%s,%s\n"
		fmt.Printf(template, title, item)
	}
	println()
}
