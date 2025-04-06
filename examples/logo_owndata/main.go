// Copyright the geany contributors.
// SPDX-License-Identifier: MPL-2.0

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
