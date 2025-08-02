// Copyright the geany contributors.
// SPDX-License-Identifier: MPL-2.0

package geany_test

import (
	"errors"
	"os"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AlphaOne1/geany"
)

func TestPrintLogo(t *testing.T) {
	tempFile, fErr := os.CreateTemp(t.TempDir(), "testPrintLogo")

	require.NoError(t, fErr)

	defer func() { assert.NoError(t, os.Remove(tempFile.Name())) }()

	save := os.Stdout
	os.Stdout = tempFile

	printErr := geany.PrintLogo(
		"Logo {{ .Geany.GoVersion }}",
		nil)

	os.Stdout = save

	require.NoError(t, printErr, "logo printing error")

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")

	target, targetErr := os.ReadFile(tempFile.Name())
	require.NoError(t, targetErr)

	assert.Equal(t, string(target), "Logo "+buildInfo.GoVersion+"\n", "Logo does not contain go version")
}

func TestPrintSimple(t *testing.T) {
	tempFile, fErr := os.CreateTemp(t.TempDir(), "testPrintLogo")

	require.NoError(t, fErr)

	defer func() { assert.NoError(t, os.Remove(tempFile.Name())) }()

	save := os.Stdout
	os.Stdout = tempFile

	printErr := geany.PrintSimple(nil)
	os.Stdout = save

	require.NoError(t, printErr, "simple writer printing error")

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")

	target, targetErr := os.ReadFile(tempFile.Name())
	require.NoError(t, targetErr)

	assert.Contains(t, string(target), `"GoVersion": "`+buildInfo.GoVersion+`"`)
	assert.Contains(t, string(target), `"VcsModified": ""`)
	assert.Contains(t, string(target), `"VcsRevision": "unknown"`)
	assert.Contains(t, string(target), `"VcsTime": "unknown"`)
}

type BrokenNIO struct {
	cnt  int
	From int
	To   int
}

func (b *BrokenNIO) Read(in []byte) (n int, err error) {
	if (b.cnt >= b.From || b.From == 0) &&
		(b.cnt < b.To || b.To == 0) {

		b.cnt++

		return 0, errors.New("broken reader")
	}
	b.cnt++

	return len(in), nil
}

func (b *BrokenNIO) Write(in []byte) (n int, err error) {
	if (b.cnt >= b.From || b.From == 0) &&
		(b.cnt < b.To || b.To == 0) {

		b.cnt++

		return 0, errors.New("broken writer")
	}
	b.cnt++

	return len(in), nil
}

func TestBrokenLogoWriter(t *testing.T) {
	t.Parallel()

	target := BrokenNIO{}

	err := geany.PrintLogoWriter(
		&target,
		"Logo {{ .Geany.GoVersion }}",
		nil)

	require.Error(t, err, "simple writer printing error")
	require.Equal(t,
		"broken writer\ncould not write program name: broken writer",
		err.Error(),
		"not two broken writer errors")
}

func TestBrokenLogoWriterFallback(t *testing.T) {
	t.Parallel()

	target := BrokenNIO{To: 1}

	err := geany.PrintLogoWriter(
		&target,
		"Logo {{ .Geany.GoVersion }}",
		nil)

	require.Error(t, err, "simple writer printing error")
	require.Equal(t, "broken writer", err.Error(), "just one broken writer error")
}

func TestBrokenLogoWriterAtEnd(t *testing.T) {
	t.Parallel()

	target := BrokenNIO{From: 2}

	err := geany.PrintLogoWriter(
		&target,
		"Logo {{ .Geany.GoVersion }}",
		nil)

	require.Error(t, err, "simple writer printing error")
	assert.Equal(t, "could not write final newline: broken writer", err.Error(), "just one broken writer error")
}

func TestBrokenSimpleWriter(t *testing.T) {
	t.Parallel()

	target := BrokenNIO{}

	require.Error(t,
		geany.PrintSimpleWriter(
			&target,
			nil),
		"simple writer no printing error")
}

func TestBrokenSimpleWriterUserData(t *testing.T) {
	t.Parallel()

	target := BrokenNIO{From: 1}

	require.Error(t,
		geany.PrintSimpleWriter(
			&target,
			nil),
		"simple writer no printing error")
}

func TestBrokenSimpleWriterFinal(t *testing.T) {
	t.Parallel()

	target := BrokenNIO{From: 2}

	require.Error(t,
		geany.PrintSimpleWriter(
			&target,
			nil),
		"simple writer no printing error")
}

func TestBrokenLogo(t *testing.T) {
	t.Parallel()

	require.Error(t,
		geany.PrintLogo("{{ .Geany }", nil),
		"broken logo does produce error")
}

func TestLogoWriterNil(t *testing.T) {
	t.Parallel()

	require.ErrorIs(t,
		geany.PrintLogoWriter(nil, "", nil),
		geany.ErrWriterNil,
		"nil writer does not produce error")

	require.ErrorIs(t,
		geany.PrintSimpleWriter(nil, nil),
		geany.ErrWriterNil,
		"nil writer does not produce error")
}
