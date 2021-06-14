// Copyright (c) 2021 Wireleap

package restartcmd

import (
	"flag"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
)

func Cmd(arg0 string, start func(fsdir.T), stop func(fsdir.T)) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("restart", flag.ExitOnError),
		Desc:    fmt.Sprintf("Restart %s daemon", arg0),
		Run: func(fm fsdir.T) {
			var (
				pid int
				err = fm.Get(&pid, arg0+".pid")
			)
			if err == nil {
				stop(fm)

				i := 0
				for ; i < 10; i++ {
					err = syscall.Kill(pid, 0)

					if err != nil {
						break
					}

					time.Sleep(500 * time.Millisecond)
				}

				if i == 10 {
					log.Fatalf("timed out waiting for %s to stop", arg0)
				}
			}
			start(fm)
		},
	}
}
