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
