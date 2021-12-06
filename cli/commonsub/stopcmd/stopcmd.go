// Copyright (c) 2021 Wireleap

package stopcmd

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
	"github.com/wireleap/common/cli/process"
)

func Cmd(arg0 string) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("stop", flag.ExitOnError),
		Desc:    fmt.Sprintf("Stop %s daemon", arg0),
		Run: func(fm fsdir.T) {
			var (
				pid     int
				err     error
				pidfile = arg0 + ".pid"
			)
			if err = fm.Get(&pid, pidfile); err != nil {
				log.Fatalf(
					"could not get pid of %s from %s: %s",
					arg0, fm.Path(pidfile), err,
				)
			}
			if process.Exists(pid) {
				if err = process.Term(pid); err != nil {
					log.Fatalf("could not terminate %s pid %d: %s", arg0, pid, err)
				}
			}
			for i := 0; i < 30; i++ {
				if !process.Exists(pid) {
					log.Printf("stopped %s daemon (was pid %d)", arg0, pid)
					fm.Del(pidfile)
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
			process.Kill(pid)
			time.Sleep(100 * time.Millisecond)
			if process.Exists(pid) {
				log.Fatalf("timed out waiting for %s (pid %d) to shut down -- process still alive!", arg0, pid)
			}
			log.Printf("stopped %s daemon (was pid %d)", arg0, pid)
			fm.Del(pidfile)
		},
	}
}
