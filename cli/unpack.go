// Copyright (c) 2022 Wireleap

package cli

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/wireleap/common/cli/fsdir"
)

var osSuffix = "_" + runtime.GOOS

// Unpack go v1.16 embedded FS contents to disk.
func UnpackEmbedded(f embed.FS, fm fsdir.T, force bool) error {
	return fs.WalkDir(f, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			// shouldn't ever happen...
			return err
		}
		origp := p
		p = strings.ReplaceAll(origp, osSuffix, "")
		if d.IsDir() {
			if e := os.MkdirAll(fm.Path(p), 0755); e != nil {
				log.Printf("could not create directory %s: %s", fm.Path(p), e)
				return e
			}
			return nil
		}
		if !force {
			if _, err := os.Stat(fm.Path(p)); err == nil {
				log.Printf("file exists: %s; not overwriting", fm.Path(p))
				return nil
			}
		}
		log.Printf("unpacking embedded file %s", p)
		b, err := f.ReadFile(origp)
		if err != nil {
			// shouldn't ever happen...
			log.Printf("error reading embedded file %s: %s", p, err)
			return err
		}
		mode := fs.FileMode(0644)
		if ok, err := filepath.Match("scripts/default/*", p); err == nil && ok && filepath.Base(p) != "README" {
			// +x
			mode = 0755
		}
		if filepath.Base(p) == "wireleap_tun" {
			// +x
			mode = 0755
		}
		if err = os.WriteFile(fm.Path(p), b, mode); err != nil {
			log.Printf("error writing embedded file %s to %s: %s", origp, fm.Path(p), err)
			return err
		}
		return nil
	})
}
