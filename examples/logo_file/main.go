// SPDX-FileCopyrightText: 2025 The geany contributors.
// SPDX-License-Identifier: MPL-2.0

// Package main contains the geany example with read logo from a file.
package main

import (
	"os"

	"github.com/AlphaOne1/geany"
)

func main() {
	logo, logoErr := os.ReadFile("logo.tmpl")

	if logoErr == nil {
		_ = geany.PrintLogo(string(logo), nil)
	}
}
