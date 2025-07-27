// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package buildinfo

import (
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfo_String(t *testing.T) {
	t.Parallel()

	zone, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	date := time.Date(2025, time.July, 25, 0, 0, 0, 0, zone)

	tests := []struct {
		name string
		info *Info
		want string
	}{
		{
			name: "WithAppVersion",
			info: &Info{
				Version:   "1.0.0",
				GoVersion: "go 1.25",
				Commit:    "qwerty",
				Time:      date,
			},
			want: "gdatum version 1.0.0 (built with go 1.25 at Fri, 25 Jul 2025 00:00:00 EDT)",
		},
		{
			name: "WithoutAppVersion",
			info: &Info{
				GoVersion: "go 1.25",
				Commit:    "qwerty",
				Time:      date,
			},
			want: "gdatum version unknown-qwerty (built with go 1.25 at Fri, 25 Jul 2025 00:00:00 EDT)",
		},
		{
			name: "Nil",
			info: nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, strings.TrimSuffix(tt.info.String(), " "+runtime.GOOS+"/"+runtime.GOARCH))
		})
	}
}
