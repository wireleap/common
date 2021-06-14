// Copyright (c) 2021 Wireleap

package superviseupgradecmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/blang/semver"
	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/commonsub/commonlib"
	"github.com/wireleap/common/cli/fsdir"
	"github.com/wireleap/common/cli/upgrade"
)

func Cmd(ctx commonlib.Context) *cli.Subcmd {
	fs := flag.NewFlagSet("supervise-upgrade", flag.ExitOnError)
	return &cli.Subcmd{
		FlagSet: fs,
		Desc:    "Supervise upgrade to a new version in a separate process (internal command)",
		Hidden:  true,
		Run: func(f fsdir.T) {
			var (
				oldbin  = f.Path(ctx.BinName + ".prev")
				curbin  = f.Path(ctx.BinName)
				newbin  = f.Path(ctx.BinName + ".next")
				oldcfg  = f.Path("config.json.prev")
				curcfg  = f.Path("config.json")
				newcfg  = f.Path("config.json.next")
				pidfile = ctx.BinName + ".pid"

				errstack []error
				err      error
			)
			defer func() {
				if len(errstack) > 0 {
					for _, e := range errstack {
						log.Printf("* %s\n", e)
					}
					log.Fatal("aborted supervised upgrade due to the above errors")
				}
			}()
			psh := func(err error) { errstack = append(errstack, err) }
			// try running pre hook
			if ctx.PreHook != nil {
				if err = ctx.PreHook(f); err != nil {
					psh(fmt.Errorf("error while running pre-upgrade hook: %s", err))
					return
				}
			}
			// try migrating
			from := ctx.NewVersion
			if fs.NArg() == 1 {
				var err error
				from, err = semver.Parse(fs.Arg(0))
				if err != nil {
					psh(fmt.Errorf("could not parse version '%s' to upgrade from: %s", fs.Arg(0), err))
					return
				}
			}
			log.Printf("running migrations...")
			if err = cli.RunChild(newbin, "migrate", from.String()); err != nil {
				psh(fmt.Errorf("migrate returned error %s", err))
				return
			}
			// remove old files
			log.Printf("removing .prev files if present...")
			for _, fn := range []string{oldcfg, oldbin} {
				os.Remove(fn)
			}
			// stop the old binary if running
			log.Printf("stopping old %s if running...", ctx.BinName)
			var pid int
			if err = f.Get(&pid, pidfile); err == nil && syscall.Kill(pid, 0) == nil {
				log.Printf("found old %s, pid %d", ctx.BinName, pid)
				if err = cli.RunChild(curbin, "stop"); err != nil {
					psh(fmt.Errorf("stopping old binary returned error %s", err))
					return
				}
			}
			// past this point, perform rollback on failure
			// if the execution got here the errstack is empty
			defer func() {
				if len(errstack) > 0 {
					log.Printf("handling upgrade failure...")
					// rollback!
					log.Printf("calling new binary rollback...")
					if err = cli.RunChild(newbin, "rollback"); err != nil {
						psh(fmt.Errorf("new binary rollback FAILED: %s", err))
						log.Printf("calling old binary rollback...")
						if err = cli.RunChild(oldbin, "rollback"); err != nil {
							psh(fmt.Errorf("old binary rollback FAILED: %s", err))
						}
					}
					// write this version to skip
					upgrade.NewConfig(f, curbin, false).SkipVersion(ctx.NewVersion)
				}
			}()
			// run the new binary
			log.Printf("running new %s...", ctx.BinName)
			if err = cli.RunChild(newbin, "start"); err != nil {
				psh(fmt.Errorf("'%s start' returned error %s", newbin, err))
				return
			}
			// actually replace binary with fallbacks
			log.Printf("replacing %s binary...", ctx.BinName)
			if err = os.Rename(curbin, oldbin); err != nil {
				psh(fmt.Errorf("renaming %s -> %s failed: %w", curbin, oldbin, err))
				return
			}
			if err = os.Rename(newbin, curbin); err != nil {
				psh(fmt.Errorf("renaming %s -> %s failed: %s", newbin, curbin, err))
				return
			}
			// actually replace config file with fallbacks
			// depends on config.json.next being created by migrate cmd
			log.Printf("replacing config file...")
			if err = os.Rename(curcfg, oldcfg); err != nil {
				psh(fmt.Errorf("renaming %s -> %s failed: %w", curcfg, oldcfg, err))
				return
			}
			if err = os.Rename(newcfg, curcfg); err != nil {
				psh(fmt.Errorf("renaming %s -> %s failed: %w", newcfg, curcfg, err))
				return
			}
			// try running post hook
			if ctx.PostHook != nil {
				if err = ctx.PostHook(f); err != nil {
					psh(fmt.Errorf("error while running post-upgrade hook: %s", err))
					return
				}
			}
		},
	}
}
