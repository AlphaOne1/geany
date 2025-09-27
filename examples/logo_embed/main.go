// SPDX-FileCopyrightText: 2025 The geany contributors.
// SPDX-License-Identifier: MPL-2.0

// Package main contains the geany example with embedded logo.
package main

import (
	_ "embed"

	"github.com/AlphaOne1/geany"
)

//go:embed logo.tmpl
var logo string

func main() {
	_ = geany.PrintLogo(logo, nil)
}
