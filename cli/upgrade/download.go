// Copyright (c) 2022 Wireleap

package upgrade

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/blang/semver"
)

const (
	releaseURL   = "https://github.com/wireleap/%s/releases/download/v%s/%s_%s-amd64%s"
	changelogURL = "https://raw.githubusercontent.com/wireleap/%s/master/changelogs/%s.md"
)

func basename(f string) string { return strings.TrimSuffix(f, filepath.Ext(f)) }

func maybeExe() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	} else {
		return ""
	}
}

func repoName(component string) string {
	switch basename(component) {
	case "wireleap":
		return "client"
	case "wireleap-relay":
		return "relay"
	case "wireleap-dir":
		return "dir"
	case "wireleap-auth":
		return "auth"
	case "wireleap-contract":
		return "contract"
	default:
		return basename(component)
	}
}

func (u *Config) BinaryURL(ver semver.Version) string {
	return fmt.Sprintf(releaseURL, repoName(u.binfile), ver.String(), basename(u.binfile), runtime.GOOS, maybeExe())
}

func (u *Config) HashURL(ver semver.Version) string {
	return fmt.Sprintf(releaseURL, repoName(u.binfile), ver.String(), basename(u.binfile), runtime.GOOS, maybeExe()+".hash")
}

func (u *Config) ChangelogURL(ver semver.Version) string {
	return fmt.Sprintf(changelogURL, repoName(u.binfile), ver.String())
}
