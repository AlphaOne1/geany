// Copyright the SonicWeb contributors.
// SPDX-License-Identifier: MPL-2.0

package geany

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func TestPrepareLogoDataNormal(t *testing.T) {
	logoData := prepareLogoData(nil)

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")

	assert.Nil(t, logoData.Values)
	assert.Empty(t, logoData.Geany.VcsModified, "", "Geany.VcsModified is not empty")
	assert.Equal(t, logoData.Geany.VcsTime, "unknown", "Geany.VcsTime is not unknown")
	assert.Equal(t, logoData.Geany.VcsRevision, "unknown", "Geany.VcsRevision is not unknown")
	assert.Equal(t, logoData.Geany.GoVersion, buildInfo.GoVersion, fmt.Sprintf("Geany.GoVersion is not %s", buildInfo.GoVersion))
}

func TestPrepareLogoDataMocked(t *testing.T) {
	old := getBuildInfo
	getBuildInfo = func() (*debug.BuildInfo, bool) {
		result := debug.BuildInfo{
			GoVersion: "go1.0",
			Path:      Must(os.Getwd()),
			Main: debug.Module{
				Path:    Must(os.Getwd()),
				Version: "1.0",
				Sum:     "deadbeafcafe",
				Replace: nil,
			},
			Deps: nil,
			Settings: []debug.BuildSetting{
				{
					Key:   "vcs.modified",
					Value: "true",
				},
				{
					Key:   "vcs.revision",
					Value: "becafe",
				},
				{
					Key:   "vcs.time",
					Value: "2006-01-02 15:04:05",
				},
			},
		}

		return &result, true
	}

	logoData := prepareLogoData(nil)

	getBuildInfo = old

	assert.Nil(t, logoData.Values)
	assert.Equal(t, logoData.Geany.VcsModified, "*", "Geany.VcsModified is not *")
	assert.Equal(t, logoData.Geany.VcsTime, "2006-01-02 15:04:05", "Geany.VcsTime is not desired timestamp")
	assert.Equal(t, logoData.Geany.VcsRevision, "becafe", "Geany.VcsRevision is not 'becafe'")
	assert.Equal(t, logoData.Geany.GoVersion, "go1.0", "Geany.GoVersion is not go1.0")
}

func TestPrintLogo(t *testing.T) {
	f, fErr := os.CreateTemp("", "testPrintLogo")

	assert.NoError(t, fErr)

	defer func() { assert.NoError(t, os.Remove(f.Name())) }()

	save := os.Stdout
	os.Stdout = f

	printErr := PrintLogo(
		"Logo {{ .Geany.GoVersion }}",
		nil)

	os.Stdout = save

	assert.NoError(t, printErr, "logo printing error")

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")

	target, targetErr := os.ReadFile(f.Name())
	assert.NoError(t, targetErr)

	assert.Equal(t, string(target), "Logo "+buildInfo.GoVersion+"\n", "Logo does not contain go version")
}

func TestPrintSimple(t *testing.T) {
	f, fErr := os.CreateTemp("", "testPrintLogo")

	assert.NoError(t, fErr)

	defer func() { assert.NoError(t, os.Remove(f.Name())) }()

	save := os.Stdout
	os.Stdout = f

	printErr := PrintSimple(nil)
	os.Stdout = save

	assert.NoError(t, printErr, "simple writer printing error")

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")

	target, targetErr := os.ReadFile(f.Name())
	assert.NoError(t, targetErr)

	assert.Contains(t, string(target), `"GoVersion":"`+buildInfo.GoVersion+`"`)
	assert.Contains(t, string(target), `"VcsModified":""`)
	assert.Contains(t, string(target), `"VcsRevision":"unknown"`)
	assert.Contains(t, string(target), `"VcsTime":"unknown"`)
}

type BrokenNIO struct {
	cnt int
	N   int
}

func (b *BrokenNIO) Read(in []byte) (n int, err error) {
	if b.cnt < b.N || b.N == 0 {
		b.cnt++
		return 0, errors.New("broken reader")
	}

	return len(in), nil
}

func (b *BrokenNIO) Write(in []byte) (n int, err error) {
	if b.cnt < b.N || b.N == 0 {
		b.cnt++
		return 0, errors.New("broken writer")
	}

	return len(in), nil
}

func TestBrokenLogoWriter(t *testing.T) {
	target := BrokenNIO{}

	err := PrintLogoWriter(
		&target,
		"Logo {{ .Geany.GoVersion }}",
		nil)

	assert.Error(t, err, "simple writer printing error")
	assert.Errorf(t, err, "broken writer")
	assert.Equal(t, err.Error(), "broken writer\nbroken writer", "not two broken writer errors")
}

func TestBrokenLogoWriterFallback(t *testing.T) {
	target := BrokenNIO{N: 1}

	err := PrintLogoWriter(
		&target,
		"Logo {{ .Geany.GoVersion }}",
		nil)

	assert.Error(t, err, "simple writer printing error")
	assert.Errorf(t, err, "broken writer")
	assert.Equal(t, err.Error(), "broken writer", "just one broken writer error")
}

func TestBrokenSimpleWriter(t *testing.T) {
	target := BrokenNIO{}

	assert.Error(t,
		PrintSimpleWriter(
			&target,
			nil),
		"simple writer no printing error")
}

func TestBrokenLogo(t *testing.T) {
	assert.Panics(t,
		func() {
			_ = PrintLogo("{{ .Geany }", nil)
		},
		"broken logo does not panic")
}
