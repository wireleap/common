// Copyright (c) 2021 Wireleap

package upgrade

import (
	"fmt"
	"runtime"

	"github.com/blang/semver"
)

const (
	releaseURL   = "https://github.com/wireleap/%s/releases/download/v%s/%s_%s-amd64%s"
	changelogURL = "https://raw.githubusercontent.com/wireleap/%s/master/changelogs/%s.md"
)

func repoName(component string) string {
	switch component {
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
		return component
	}
}

func (u *Config) BinaryURL(ver semver.Version) string {
	return fmt.Sprintf(releaseURL, repoName(u.binfile), ver.String(), u.binfile, runtime.GOOS, func() string {
		if runtime.GOOS == "windows" {
			return ".exe"
		} else {
			return ""
		}
	}())
}

func (u *Config) HashURL(ver semver.Version) string {
	return fmt.Sprintf(releaseURL, repoName(u.binfile), ver.String(), u.binfile, runtime.GOOS, ".hash")
}

func (u *Config) ChangelogURL(ver semver.Version) string {
	return fmt.Sprintf(changelogURL, repoName(u.binfile), ver.String())
}
