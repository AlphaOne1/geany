// SPDX-FileCopyrightText: 2025 The geany contributors.
// SPDX-License-Identifier: MPL-2.0

package geany

import (
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
	t.Parallel()

	logoData := prepareLogoData(nil)

	buildInfo, ok := debug.ReadBuildInfo()

	assert.True(t, ok, "build info not found in debug.ReadBuildInfo")

	assert.Nil(t, logoData.Values)
	assert.Empty(t, logoData.Geany.VcsModified, "Geany.VcsModified is not empty")
	assert.Equal(t, "unknown", logoData.Geany.VcsTime, "Geany.VcsTime is not unknown")
	assert.Equal(t, "unknown", logoData.Geany.VcsRevision, "Geany.VcsRevision is not unknown")
	assert.Equal(t, buildInfo.GoVersion, logoData.Geany.GoVersion, "Geany.GoVersion is not "+buildInfo.GoVersion)
}

func TestPrepareLogoDataMocked(t *testing.T) {
	t.Parallel()

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
	assert.Equal(t, "*", logoData.Geany.VcsModified, "Geany.VcsModified is not *")
	assert.Equal(t, "2006-01-02 15:04:05", logoData.Geany.VcsTime, "Geany.VcsTime is not desired timestamp")
	assert.Equal(t, "becafe", logoData.Geany.VcsRevision, "Geany.VcsRevision is not 'becafe'")
	assert.Equal(t, "go1.0", logoData.Geany.GoVersion, "Geany.GoVersion is not go1.0")
}
