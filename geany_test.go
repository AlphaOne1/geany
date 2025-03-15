package main

import (
	"bytes"
	"errors"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareLogoData(t *testing.T) {
	logoData := prepareLogoData(nil)

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")

	assert.Nil(t, logoData.Values)
	assert.Empty(t, logoData.Geany.VcsModified, "", "Geany.VcsModified is not empty")
	assert.Equal(t, logoData.Geany.VcsTime, "unknown", "Geany.VcsTime is not unknown")
	assert.Equal(t, logoData.Geany.VcsRevision, "unknown", "Geany.VcsRevision is not unknown")
	assert.Equal(t, logoData.Geany.GoVersion, buildInfo.GoVersion, fmt.Sprintf("Geany.GoVersion is not %s", buildInfo.GoVersion))
}

func TestPrintLogoWriter(t *testing.T) {
	target := bytes.Buffer{}

	assert.Nil(t, PrintLogoWriter(
		&target,
		"Logo {{ .Geany.GoVersion }}",
		nil), "logo printing error")

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")
	assert.Equal(t, target.String(), "Logo "+buildInfo.GoVersion, "Logo does not contain go version")
}

func TestPrintSimpleWriter(t *testing.T) {
	target := bytes.Buffer{}

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")
	assert.Nil(t, PrintSimpleWriter(&target, nil), "simple writer printing error")

	assert.Contains(t, target.String(), `"GoVersion":"`+buildInfo.GoVersion+`"`)
	assert.Contains(t, target.String(), `"VcsModified":""`)
	assert.Contains(t, target.String(), `"VcsRevision":"unknown"`)
	assert.Contains(t, target.String(), `"VcsTime":"unknown"`)
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

func TestBrokenLogoWiter(t *testing.T) {
	target := BrokenNIO{}

	assert.NotNil(t,
		PrintLogoWriter(
			&target,
			"Logo {{ .Geany.GoVersion }}",
			nil),
		"simple writer no printing error")
}

func TestBrokenSimpleWiter(t *testing.T) {
	target := BrokenNIO{}

	assert.NotNil(t,
		PrintSimpleWriter(
			&target,
			nil),
		"simple writer no printing error")
}
