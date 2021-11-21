// Copyright (c) 2021 Wireleap

package reloadcmd

import (
	"flag"
	"fmt"

	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
)

func Cmd(arg0 string) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("reload", flag.ExitOnError),
		Desc:    fmt.Sprintf("Reload %s daemon configuration (NOT IMPLEMENTED YET)", arg0),
		Run: func(fm fsdir.T) {
			fmt.Printf("`reload` is not implemented on Windows yet. Please use `restart` instead.")
		},
	}
}
