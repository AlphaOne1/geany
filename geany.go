// Copyright the geany contributors.
// SPDX-License-Identifier: MPL-2.0

// Package geany contains the logo printing functionality.
package geany

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"text/template"
)

// geanyData contains the data that geany provies on itself.
type geanyData struct {
	VcsRevision string
	VcsTime     string
	VcsModified string
	GoVersion   string
}

// logoData is the structure given as data to the templating engine.
// It contains the geany provided build information and the user provided data.
type logoData struct {
	Geany  geanyData
	Values any
}

// getBuildInfo is a global variable to mock it easily in tests.
var getBuildInfo = debug.ReadBuildInfo //nolint:gochecknoglobals

// prepareLogoData collects the build information and the user-provided data
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

	if buildInfo, ok := getBuildInfo(); ok {
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

// PrintSimpleWriter outputs just the name of the program, the build information
// and, in case, the user given data to a user provided io.Writer.
func PrintSimpleWriter(writer io.Writer, values any) error {
	revData := prepareLogoData(values)

	// normally we have the program's name given as the first argument
	if len(os.Args) > 0 {
		fmt.Printf("%s\n", os.Args[0])
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")

	// we suppress the linter here, as we cannot guarantee for users data.
	err := encoder.Encode(revData) //nolint:musttag

	if err == nil {
		_, err = fmt.Fprintln(writer)
	}

	if err != nil {
		return fmt.Errorf("could not write: %w", err)
	}

	return nil
}

// PrintSimple is a convenience wrapper around PrintSimpleWriter,
// as logos are normally printed to standard output.
func PrintSimple(values any) error {
	return PrintSimpleWriter(os.Stdout, values)
}

// PrintLogoWriter takes a text/template string as its parameter and renders it to be the logo.
// It offers the following data for the template:
//   - VcsRevision
//   - VcsTime
//   - VcsModified
//   - GoVersion
//
// these can be referenced in the template, e.g., using {{ .VcsRevision }}.
// An additional custom value can be accessed via the Values field. Its type must match the way
// that it is accessed in the logo.
func PrintLogoWriter(writer io.Writer, tmpl string, values any) error {
	revData := prepareLogoData(values)

	logo := template.New("logo")
	template.Must(logo.Parse(tmpl))

	if err := logo.Execute(writer, revData); err != nil {
		return errors.Join(err, PrintSimpleWriter(writer, values))
	}

	_, err := fmt.Fprintln(writer)

	if err != nil {
		return fmt.Errorf("could not write: %w", err)
	}

	return nil
}

// PrintLogo is a convenience wrapper around PrintLogoWriter,
// as logos are normally printed to standard output.
func PrintLogo(tmpl string, values any) error {
	return PrintLogoWriter(os.Stdout, tmpl, values)
}
