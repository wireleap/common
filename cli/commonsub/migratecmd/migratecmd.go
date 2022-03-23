// Copyright (c) 2022 Wireleap

package migratecmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/blang/semver"
	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
	"github.com/wireleap/common/cli/upgrade"
)

func Cmd(binname string, ms0 []*upgrade.Migration, v semver.Version) *cli.Subcmd {
	r := &cli.Subcmd{
		Hidden:  true,
		FlagSet: flag.NewFlagSet("migrate", flag.ExitOnError),
		Desc:    fmt.Sprintf("Migrate %s files (internal command)", binname),
	}
	r.Run = func(f fsdir.T) {
		if r.FlagSet.NArg() != 1 {
			log.Fatalf("which version to migrate from? usage: `%s migrate 1.2.3`", binname)
		}
		v0, err := semver.Parse(r.FlagSet.Arg(0))
		if err != nil {
			log.Fatalf("could not parse supplied version '%s': %s", r.FlagSet.Arg(0), err)
		}
		b, err := ioutil.ReadFile(f.Path("config.json"))
		if err != nil {
			log.Fatalf("could not read old config file: %s", err)
		}
		if err = ioutil.WriteFile(f.Path("config.json.next"), b, 0644); err != nil {
			log.Fatalf("could not copy old config to new config file: %s", err)
		}
		if ms0 != nil {
			ms := []*upgrade.Migration{}
			for _, m := range ms0 {
				if m.Version.GT(v0) && m.Version.LTE(v) {
					ms = append(ms, m)
				}
			}
			sort.Slice(ms, func(i, j int) bool { return ms[i].Version.LT(ms[j].Version) })
			for _, m := range ms {
				log.Printf(
					"running %s migration '%s' for %s %s...",
					m.Version, m.Name, binname, v,
				)
				if err := m.TryApply(f); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	return r
}
