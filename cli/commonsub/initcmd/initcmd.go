// Copyright (c) 2021 Wireleap

package initcmd

import (
	"crypto/ed25519"
	"embed"
	"flag"
	"fmt"
	"log"

	"github.com/wireleap/common/api/jsonb"
	"github.com/wireleap/common/cli"
	"github.com/wireleap/common/cli/fsdir"
)

const (
	Seed = "key.seed"
	Pub  = "key.pub"
)

func Cmd(arg0 string, steps ...func(fsdir.T) error) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("init", flag.ExitOnError),
		Desc:    "Initialize %s files",
		Run: func(fm fsdir.T) {
			for _, s := range steps {
				if err := s(fm); err != nil {
					log.Fatal("error while initializing:", err)
				}
			}
		},
	}
}

func KeypairStep(f fsdir.T) error {
	pk, sk, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}
	log.Printf("writing ed25519 private key seed to %s", f.Path(Seed))
	if err = f.Set(jsonb.B(sk.Seed()), Seed); err != nil {
		return err
	}
	if err = f.Chmod(0600, Seed); err != nil {
		return err
	}
	log.Printf("writing ed25519 public key to %s", f.Path(Pub))
	if err = f.Set(jsonb.PK(pk), Pub); err != nil {
		return err
	}
	return nil
}

func UnpackStep(fs embed.FS) func(fsdir.T) error {
	return func(f fsdir.T) error {
		if err := cli.UnpackEmbedded(fs, f, false); err != nil {
			return fmt.Errorf("error unpacking embedded files: %s", err)
		}
		return nil
	}
}
