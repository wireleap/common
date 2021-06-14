// Copyright (c) 2021 Wireleap

package versioncmd

import (
	"flag"
	"fmt"

	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
)

func Cmd(vstring string) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("version", flag.ExitOnError),
		Desc:    "Show version and exit",
		Run:     func(_ fsdir.T) { fmt.Println(vstring) },
	}
}
