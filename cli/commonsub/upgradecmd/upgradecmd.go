// Copyright (c) 2022 Wireleap

package upgradecmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/blang/semver"
	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
	"github.com/wireleap/common/cli/upgrade"
)

func Cmd(arg0 string, ex upgrade.Executor, v0 semver.Version, v1f func(fsdir.T) (semver.Version, error)) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("upgrade", flag.ExitOnError),
		Desc:    fmt.Sprintf("Upgrade %s to the latest version per directory", arg0),
		Run: func(fm fsdir.T) {
			if !upgrade.Supported {
				fmt.Printf("Your %s binary does not support upgrades (manual build?), aborting.\n", arg0)
				os.Exit(1)
			}
			die := func(e error) {
				if e != nil {
					log.Fatalf("error: %s", e)
				}
			}
			v1, err := v1f(fm)
			die(err)
			// check if latest <= current
			if v1.LE(v0) {
				fmt.Printf(
					"Your %s is up to date (version %s, %s available), nothing to do.\n",
					arg0, v0, v1,
				)
				return
			}
			// initialize upgrade config
			u := upgrade.NewConfig(fm, arg0, true)
			// test
			// download & present changelog
			chglog, err := u.GetChangelog(v1)
			if err != nil {
				chglog = fmt.Sprintf("-- error getting changelog: %s --", err)
			}
			fmt.Printf("Changelog for version %s:\n%s\n", v1, chglog)
			if sv := u.SkippedVersion(); sv != nil && sv.EQ(v1) {
				fmt.Printf("-- NOTE: upgrading to version %s seems to have failed the last time --\n", v1)
			}
			// confirm intention
			if !upgrade.Confirm(fmt.Sprintf("Proceed with upgrade from %s to %s?", v0, v1)) {
				fmt.Println("OK, aborting.")
				u.Cleanup()
				os.Exit(1)
			}
			// actually upgrade
			die(u.Upgrade(ex, v0, v1))
			fmt.Printf("%s successfully upgraded to version %s (was %s).\n", arg0, v1, v0)
			// clean up only in case of success
			u.Cleanup()
			// warn about manual rollback
			fmt.Printf(
				"\nNOTE: The old (%s) binary can be restored:\ncd %s && ./%s rollback\n",
				v0, fm.Path(), arg0,
			)
		},
	}
}
