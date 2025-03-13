// Copyright the SonicWeb contributors.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"text/template"
)

// geanyData contains the data that geany provies on itself
type geanyData struct {
	VcsRevision string
	VcsTime     string
	VcsModified string
	GoVersion   string
}

// logoData is the structure that is given as data to the templating engine.
// It contains the geany provided build information and the user provided data.
type logoData struct {
	Geany  geanyData
	Values any
}

// prepareLogoData collects the build information and the user provided data
// into a logoData structure.
func prepareLogoData(values any) logoData {
	result := logoData{
		Geany: geanyData{
			VcsRevision: "unknown",
			VcsTime:     "unknown",
			VcsModified: "",
			GoVersion:   "unknown",
		},
		Values: values,
	}

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		result.Geany.GoVersion = buildInfo.GoVersion

		for _, s := range buildInfo.Settings {
			switch s.Key {
			case "vcs.revision":
				result.Geany.VcsRevision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					result.Geany.VcsModified = "*"
				}
			case "vcs.time":
				result.Geany.VcsTime = s.Value
			}
		}
	}

	return result
}

// PrintSimpleWriter outputs just the name of the programm, the build information
// and, in case, the user given data to a user provided io.Writer.
func PrintSimpleWriter(w io.Writer, values any) error {
	revData := prepareLogoData(values)

	// normally we have the programs name given as first argument
	if len(os.Args) > 0 {
		fmt.Printf("%s\n", os.Args[0])
	}

	return json.NewEncoder(w).Encode(revData)
}

// PrintSimple is a convenience wrapper around PrintSimpleWriter,
// as logos are normally printed to standard output.
func PrintSimple(values any) error {
	return PrintSimpleWriter(os.Stdout, values)
}

// PrintLogoWriter takes a text/template string as parameter and renders it to be the logo. It offers the
// following data for the template:
//   - VcsRevision
//   - VcsTime
//   - VcsModified
//   - GoVersion
//
// these can be referenced in the template, e.g. using {{ .VcsRevision }}.
// An additional custom value can be accessed via the Values field. Its type must match the way
// that it is accessed in the logo.
func PrintLogoWriter(w io.Writer, tmpl string, values any) error {
	revData := prepareLogoData(values)

	logo := template.New("logo")
	template.Must(logo.Parse(tmpl))

	if err := logo.Execute(w, revData); err != nil {
		return errors.Join(err, PrintSimpleWriter(w, values))
	}

	fmt.Println()
	return nil
}

// PrintLogo is a convenience wrapper around PrintLogoWriter,
// as logos are normally printed to standard output.
func PrintLogo(tmpl string, values any) error {
	return PrintLogoWriter(os.Stdout, tmpl, values)
}
