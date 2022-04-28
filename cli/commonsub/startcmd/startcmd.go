// Copyright (c) 2022 Wireleap

package startcmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/tabwriter"
	"time"

	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
	"github.com/wireleap/common/cli/process"
)

func Cmd(arg0 string, do func(fsdir.T)) *cli.Subcmd {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	fg := fs.Bool("fg", false, "Run in foreground, don't detach")

	r := &cli.Subcmd{
		FlagSet: fs,
		Desc:    fmt.Sprintf("Start %s daemon", arg0),
		Run: func(fm fsdir.T) {
			var err error

			if *fg == false {
				var pid int
				if err = fm.Get(&pid, arg0+".pid"); err == nil {
					if process.Exists(pid) {
						log.Fatalf("%s daemon is already running!", arg0)
					}
				}

				binary, err := exec.LookPath(os.Args[0])
				if err != nil {
					log.Fatalf("could not find own binary path: %s", err)
				}

				logpath := fm.Path(arg0 + ".log")
				logfile, err := os.OpenFile(logpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

				if err != nil {
					log.Fatalf("could not open logfile %s: %s", logpath, err)
				}
				defer logfile.Close()

				cmd := exec.Cmd{
					Path:   binary,
					Args:   []string{binary, "start", "--fg"},
					Stdout: logfile,
					Stderr: logfile,
					Env:    append(os.Environ(), "WIRELEAP_BACKGROUND=1"),
				}
				if err = cmd.Start(); err != nil {
					log.Fatalf("could not spawn background %s process: %s", arg0, err)
				}
				log.Printf(
					"starting %s with pid %d, writing to %s...",
					arg0, cmd.Process.Pid, logpath,
				)
				// wait for 2s and see if it's still alive
				e := make(chan error)
				go func() { e <- cmd.Wait() }()
				select {
				case <-e:
					log.Printf("%s is not running, %s follows:", arg0, logpath)
					b, err := ioutil.ReadFile(logpath)
					if err != nil {
						log.Fatalf("could not get %s contents!", logpath)
					}
					os.Stdout.Write(b)
					os.Exit(1)
				case <-time.NewTimer(time.Second * 2).C:
					log.Printf(
						"successfully spawned %s with pid %d, writing to %s",
						arg0, cmd.Process.Pid, logpath,
					)
				}
				return
			}

			err = fm.Set(os.Getpid(), arg0+".pid")

			if err != nil {
				log.Fatalf("could not write pid: %s", err)
			}

			do(fm)
		},
	}
	r.Writer = tabwriter.NewWriter(r.FlagSet.Output(), 6, 8, 1, ' ', 0)
	return r
}
