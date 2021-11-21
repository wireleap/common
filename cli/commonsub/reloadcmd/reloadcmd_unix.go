// +build linux darwin
// Copyright (c) 2021 Wireleap

package reloadcmd

import (
	"flag"
	"fmt"
	"log"
	"syscall"

	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
)

func Cmd(arg0 string) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("reload", flag.ExitOnError),
		Desc:    fmt.Sprintf("Reload %s daemon configuration", arg0),
		Run: func(fm fsdir.T) {
			var (
				pid int
				err = fm.Get(&pid, arg0+".pid")
			)

			if err != nil {
				log.Fatalf(
					"could not get pid of %s from %s: %s",
					arg0,
					fm.Path(arg0+".pid"),
					err,
				)
			}

			err = syscall.Kill(pid, syscall.SIGUSR1)

			if err != nil {
				log.Fatalf(
					"could not reload %s pid %d: %s",
					arg0,
					pid,
					err,
				)
			}

			log.Printf("reloaded %s daemon (pid %d)", arg0, pid)
		},
	}
}
