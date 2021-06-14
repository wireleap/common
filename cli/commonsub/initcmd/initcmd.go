// Copyright (c) 2021 Wireleap

package initcmd

import (
	"crypto/ed25519"
	"flag"
	"log"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
)

const (
	Seed = "key.seed"
	Pub  = "key.pub"
)

var Cmd = &cli.Subcmd{
	FlagSet: flag.NewFlagSet("init", flag.ExitOnError),
	Desc:    "Generate ed25519 keypair (key.seed, key.pub) and exit",
	Run: func(fm fsdir.T) {
		pk, sk, err := ed25519.GenerateKey(nil)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("writing ed25519 private key seed to %s", fm.Path(Seed))
		err = fm.Set(jsonb.B(sk.Seed()), Seed)

		if err != nil {
			log.Fatal(err)
		}

		err = fm.Chmod(0600, Seed)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("writing ed25519 public key to %s", fm.Path(Pub))
		err = fm.Set(jsonb.PK(pk), Pub)

		if err != nil {
			log.Fatal(err)
		}
	},
}
