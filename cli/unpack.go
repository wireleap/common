// Copyright (c) 2021 Wireleap

package cli

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/wireleap/common/cli/fsdir"
)

// Unpack go v1.16 embedded FS contents to disk.
func UnpackEmbedded(f embed.FS, fm fsdir.T, force bool) error {
	return fs.WalkDir(f, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			// shouldn't ever happen...
			return err
		}
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
		b, err := f.ReadFile(p)
		if err != nil {
			// shouldn't ever happen...
			log.Printf("error reading embedded file %s: %s", p, err)
			return err
		}
		mode := fs.FileMode(0755)
		switch path.Ext(p) {
		case "so", "html", "js":
			// no need to +x
			mode = 0644
		}
		switch path.Base(p) {
		case "README", "DISCLAIMER":
			// no need to +x
			mode = 0644
		}
		if err = os.WriteFile(fm.Path(p), b, mode); err != nil {
			log.Printf("error writing embedded file %s to %s: %s", p, fm.Path(p), err)
			return err
		}
		return nil
	})
}
