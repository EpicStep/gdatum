// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package buildinfo

import (
	"runtime/debug"
	"time"
)

const (
	projectPackage = "github.com/EpicStep/gdatum"
)

var (
	// binaryVersion is set by ldflag
	binaryVersion string

	info *Info
)

// Get Info.
func Get() *Info {
	return info
}

func initInfo() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	info = &Info{
		GoVersion: buildInfo.GoVersion,
		Version:   binaryVersion,
	}

	var isDep bool

	if binaryVersion == "" {
		if version, ok := getPackageVersion(&buildInfo.Main); ok {
			info.Version = version
		} else {
			isDep = true
			for _, m := range buildInfo.Deps {
				if v, ok := getPackageVersion(m); ok {
					info.Version = v
					break
				}
			}
		}
	}

	if !isDep {
		for _, setting := range buildInfo.Settings {
			switch setting.Key { //revive:disable:enforce-switch-style
			case "vcs.revision":
				info.Commit = setting.Value
			case "vcs.time":
				if t, err := time.Parse(time.RFC3339Nano, setting.Value); err == nil {
					info.Time = t
				}
			}
		}
	}
}

func getPackageVersion(m *debug.Module) (string, bool) {
	if m == nil || m.Path != projectPackage {
		return "", false
	}
	return m.Version, true
}

func init() {
	initInfo()
}
