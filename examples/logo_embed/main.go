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
