// SPDX-FileCopyrightText: 2026 The geany contributors.
// SPDX-License-Identifier: MPL-2.0

// Package main contains the geany example with embedded logo and custom data.
package main

import (
	_ "embed"

	"github.com/AlphaOne1/geany"
)

//go:embed logo.tmpl
var logo string

func main() {
	_ = geany.PrintLogo(logo,
		&struct {
			FeatureA bool
			FeatureB bool
			Greeting string
		}{
			FeatureA: true,
			FeatureB: false,
			Greeting: "Hi Geany!",
		})
}
