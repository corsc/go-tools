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

const (
	keyChild    = "Child Packages"
	keyStdLib   = "Std Lib Packages"
	keyExternal = "External Packages"
	keyVendored = "Vendored Packages"
)

// PrintSummary will print a summary of dependencies
func PrintSummary(in *Summary) {
	doPrintSummary(keyChild, in.child)
	doPrintSummary(keyExternal, in.external)
	doPrintSummary(keyStdLib, in.stdLib)
	doPrintSummary(keyVendored, in.vendored)
}

func doPrintSummary(title string, itemsMap map[string]*SummaryItem) {
	spacer := "----------------------------------------------------------------------------------------------------------------------------------------------------------------------------\n"

	fmt.Print(spacer)
	header := "| %s | %-161s |\n"
	fmt.Printf(header, "Deps", title)
	fmt.Print(spacer)

	sortedItems := make([]string, 0, len(itemsMap))
	for key := range itemsMap {
		sortedItems = append(sortedItems, key)
	}
	sort.Strings(sortedItems)

	for _, pkgName := range sortedItems {
		summary := itemsMap[pkgName]

		template := "| %4d | %-161s |\n"
		fmt.Printf(template, len(summary.Dependents), pkgName)
	}

	fmt.Print(spacer)
	println()
}

// PrintSummaryCSV will print a summary of dependencies to CSV
func PrintSummaryCSV(in *Summary) {
	doPrintSummaryCSV(keyChild, in.child)
	doPrintSummaryCSV(keyExternal, in.external)
	doPrintSummaryCSV(keyStdLib, in.stdLib)
	doPrintSummaryCSV(keyVendored, in.vendored)
}

func doPrintSummaryCSV(title string, itemsMap map[string]*SummaryItem) {
	fmt.Printf("** %s **\n", title)

	sortedItems := make([]string, 0, len(itemsMap))
	for key := range itemsMap {
		sortedItems = append(sortedItems, key)
	}
	sort.Strings(sortedItems)

	for _, pkgName := range sortedItems {
		summary := itemsMap[pkgName]

		template := "%d,%s\n"
		fmt.Printf(template, len(summary.Dependents), pkgName)
	}
	println()
}
