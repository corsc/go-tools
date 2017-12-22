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
	"sort"
)

// PrintDirect will output the direct deps
func PrintDirect(in *Summary) {
	printList(keyDirect, in.direct)
}

// PrintChild will output the child deps
func PrintChild(in *Summary) {
	printList(keyChild, in.child)
}

// PrintStdLib will output the std lib deps
func PrintStdLib(in *Summary) {
	printList(keyStdLib, in.stdLib)
}

// PrintExternal will output the external deps
func PrintExternal(in *Summary) {
	printList(keyExternal, in.external)
}

// PrintVendored will output the vendored deps
func PrintVendored(in *Summary) {
	printList(keyVendored, in.vendored)
}

func printList(title string, items map[string]struct{}) {
	sortedItems := make([]string, 0, len(items))
	for key := range items {
		sortedItems = append(sortedItems, key)
	}
	sort.Strings(sortedItems)

	header := "\n%-30s\n"
	fmt.Printf(header, title)
	fmt.Print("------------------------------\n")

	template := "%s\n"
	for _, item := range sortedItems {
		fmt.Printf(template, item)
	}
	println()
}
