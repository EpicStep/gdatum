// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package buildinfo

import (
	"runtime"
	"strings"
	"time"
)

// Info is the gdatum build information.
type Info struct {
	// Version is the version of the gdatum.
	Version string
	// GoVersion is the version of the Go that produced this binary.
	GoVersion string
	// Commit is the current commit hash.
	Commit string
	// Time is the time of the build.
	Time time.Time
}

func (i *Info) String() string {
	if i == nil {
		return ""
	}

	var s strings.Builder
	s.WriteString("gdatum version ")
	if v := i.Version; v != "" {
		s.WriteString(v)
	} else {
		s.WriteString("unknown")
	}

	if isDevelVersion(i.Version) && i.Commit != "" {
		s.WriteByte('-')
		s.WriteString(i.Commit)
	}

	if t, v := i.Time, i.GoVersion; v != "" || !t.IsZero() {
		s.WriteString(" (built")
		if v != "" {
			s.WriteString(" with ")
			s.WriteString(v)
		}
		if !t.IsZero() {
			s.WriteString(" at ")
			s.WriteString(t.Format(time.RFC1123))
		}
		s.WriteByte(')')
	}

	s.WriteString(" " + runtime.GOOS + "/" + runtime.GOARCH)
	return s.String()
}

func isDevelVersion(ver string) bool {
	return ver == "" || ver == "(devel)"
}
