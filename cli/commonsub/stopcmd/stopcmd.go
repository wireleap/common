// Copyright (c) 2021 Wireleap

package stopcmd

import (
	"flag"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
)

func Cmd(arg0 string) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("stop", flag.ExitOnError),
		Desc:    fmt.Sprintf("Stop %s daemon", arg0),
		Run: func(fm fsdir.T) {
			var (
				pid int
				err error
			)
			if err = fm.Get(&pid, arg0+".pid"); err != nil {
				log.Fatalf(
					"could not get pid of %s from %s: %s",
					arg0, fm.Path(arg0+".pid"), err,
				)
			}
			if err = syscall.Kill(pid, syscall.SIGTERM); err != nil {
				log.Fatalf("could not send SIGTERM to %s pid %d: %s", arg0, pid, err)
			}
			for i := 0; i < 10; i++ {
				if err = syscall.Kill(pid, 0); err != nil {
					log.Printf("stopped %s daemon (was pid %d)", arg0, pid)
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
			syscall.Kill(pid, syscall.SIGKILL)
			time.Sleep(100 * time.Millisecond)
			if err = syscall.Kill(pid, 0); err == nil {
				log.Fatalf("timed out waiting for %s (pid %d) to shut down -- process still alive!", arg0, pid)
			}
			log.Printf("stopped %s daemon (was pid %d)", arg0, pid)
		},
	}
}
