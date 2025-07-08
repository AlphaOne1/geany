// Copyright the geany contributors.
// SPDX-License-Identifier: MPL-2.0

// Package main contains the geany example with only the simple logo output.
package main

import (
	_ "embed"

	"github.com/AlphaOne1/geany"
)

func main() {
	_ = geany.PrintSimple(nil)
}
