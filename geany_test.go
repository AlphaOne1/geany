package main

import (
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
