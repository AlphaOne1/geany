package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	logoTemplateFile, logoTemplateFileErr := os.Open("testData/logo.tpl")

	assert.Nil(t, logoTemplateFileErr, "could not open template file")

	defer func() { _ = logoTemplateFile.Close() }()

	logoTemplate, readErr := io.ReadAll(logoTemplateFile)

	assert.Nil(t, readErr, "could not read logo template")

	target := bytes.Buffer{}

	assert.Nil(t, PrintLogoWriter(&target, string(logoTemplate), nil), "logo printing error")

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")
	assert.Contains(t, target.String(), buildInfo.GoVersion, "Logo does not contain go version")
}
