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

// ErrWriterNil indicates that the provided io.Writer is nil and cannot be used.
var ErrWriterNil = errors.New("writer is nil")

// geanyData contains the data that geany provides on itself.
type geanyData struct {
	VcsRevision string
	VcsTime     string
	VcsModified string
	GoVersion   string
}

// logoData is the structure given as data to the templating engine.
// It contains the geany provided build information and the user provided data.
// The user is not required to pass data, in that case nil should be passed. In any case,
// the data accessed in the logo template and provided data must match.
type logoData struct {
	Geany  geanyData // build information
	Values any       // user provided data
}

// getBuildInfo is a global variable to mock it easily in tests. Take into account that
// this variable can produce race conditions in parallel testing.
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

	if buildInfo, ok := getBuildInfo(); ok && buildInfo != nil {
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
	if writer == nil {
		return ErrWriterNil
	}

	revData := prepareLogoData(values)

	// normally we have the program's name given as the first argument
	if len(os.Args) > 0 && os.Args[0] != "" {
		if _, err := fmt.Fprintf(writer, "%s\n", os.Args[0]); err != nil {
			return fmt.Errorf("could not write program name: %w", err)
		}
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "    ")

	// we suppress the linter here, as we cannot guarantee for users data.
	if err := encoder.Encode(revData); err != nil { //nolint:musttag
		return fmt.Errorf("could not encode user data: %w", err)
	}

	if _, err := fmt.Fprintln(writer); err != nil {
		return fmt.Errorf("could not write final newline: %w", err)
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
// these can be referenced in the template, e.g., using `{{ .Geany.VcsRevision }}`.
// An additional custom value can be accessed via the Values field. Its type must match the way
// that it is accessed in the logo and is accessed using e.g. `{{ .Values.Foo }}`.
//
// The template is parsed and executed. In case of an error, the program's name and user given
// data are printed as JSON as fallback. The original error and, in case, the error of the fallback
// are returned.
func PrintLogoWriter(writer io.Writer, tmpl string, values any) error {
	if writer == nil {
		return ErrWriterNil
	}

	revData := prepareLogoData(values)

	logo := template.New("logo")

	if _, err := logo.Parse(tmpl); err != nil {
		return fmt.Errorf("could not parse template: %w", err)
	}

	if err := logo.Execute(writer, revData); err != nil {
		return errors.Join(err, PrintSimpleWriter(writer, values))
	}

	if _, err := fmt.Fprintln(writer); err != nil {
		return fmt.Errorf("could not write final newline: %w", err)
	}

	return nil
}

// PrintLogo is a convenience wrapper around PrintLogoWriter,
// as logos are normally printed to standard output.
func PrintLogo(tmpl string, values any) error {
	return PrintLogoWriter(os.Stdout, tmpl, values)
}
