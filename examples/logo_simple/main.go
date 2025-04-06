// Copyright the geany contributors.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	_ "embed"

	"github.com/AlphaOne1/geany"
)

func main() {
	_ = geany.PrintSimple(nil)
}
