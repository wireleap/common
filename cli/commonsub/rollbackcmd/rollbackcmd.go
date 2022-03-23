// Copyright (c) 2022 Wireleap

package rollbackcmd

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/commonsub/commonlib"
	"github.com/wireleap/common/cli/fsdir"
)

func Cmd(ctx commonlib.Context) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("rollback", flag.ExitOnError),
		Desc:    "Undo a partially completed upgrade",
		Run: func(f fsdir.T) {
			var (
				curbin = f.Path(ctx.BinName)
				oldbin = f.Path(ctx.BinName + ".prev")
				curcfg = f.Path("config.json")
				oldcfg = f.Path("config.json.prev")
			)
			// try running pre hook
			if ctx.PreHook != nil {
				if err := ctx.PreHook(f); err != nil {
					log.Fatalf("error while running pre-rollback hook: %s", err)
				}
			}
			// stop currently running binary if running
			// not a hard fail
			log.Printf("stopping running %s (if present)...", ctx.BinName)
			if o, err := exec.Command(curbin, "stop").CombinedOutput(); err != nil {
				log.Printf("'%s stop' returned error %s, output: %s", curbin, err, string(o))
			}
			// binary
			if _, err := os.Stat(curbin); err == nil {
				if _, err := os.Stat(oldbin); err == nil {
					log.Printf("restoring old binary...")
					if err := os.Rename(oldbin, curbin); err != nil {
						log.Fatalf("could not move %s to %s", oldbin, curbin)
					}
				}
			}
			// config file
			if _, err := os.Stat(curcfg); err == nil {
				if _, err := os.Stat(oldcfg); err == nil {
					log.Printf("restoring old config file...")
					if err := os.Rename(oldcfg, curcfg); err != nil {
						log.Fatalf("could not move %s to %s", oldcfg, curcfg)
					}
				}
			}
			// run post-rollback hook
			if ctx.PostHook != nil {
				if err := ctx.PostHook(f); err != nil {
					log.Fatalf("error while running post-rollback hook: %s", err)
				}
			}
		},
	}
}
