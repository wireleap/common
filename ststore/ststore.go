// Copyright (c) 2022 Wireleap

// The ststore package provides a concurrent in-memory sharetoken store which
// is synced to disk after modifications.
package ststore

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/wireleap/common/api/sharetoken"
	"github.com/wireleap/common/cli/fsdir"
)

const pathSeparator = string(os.PathSeparator)

type (
	st1map map[string]*sharetoken.T
	st2map map[string]map[string]*sharetoken.T
	st3map map[string]map[string]map[string]*sharetoken.T
)

// DuplicateSTError is returned if a sharetoken which was already seen is
// being added to the store.
var DuplicateSTError = errors.New("duplicate sharetoken")

// T is the type of a sharetoken store.
type T struct {
	m    fsdir.T
	mu   sync.RWMutex
	sts  st3map
	keyf KeyFunc
}

// New initializes a sharetoken store in the directory under the path given by
// the dir argument.
func New(dir string, keyf KeyFunc) (t *T, err error) {
	t = &T{keyf: keyf, sts: st3map{}}
	t.m, err = fsdir.New(dir)

	if err != nil {
		return
	}

	initPath := t.m.Path()
	initPathDepth := len(strings.Split(initPath, pathSeparator))

	err = filepath.Walk(t.m.Path(), func(path string, info os.FileInfo, err error) error {
		pathSlice := strings.Split(path, pathSeparator)
		depth := len(pathSlice)

		switch {
		case err != nil:
			return err

		case info.IsDir():
			// Limit depth to 1 level
			if depth <= initPathDepth+1 {
				switch pathSlice[depth-1] {
				// Exclude contracts with reserved naming
				case "malformed", "expired":
					return filepath.SkipDir
				default:
					return nil
				}
			}

			return filepath.SkipDir

		case depth <= initPathDepth+1:
			// Exclude files in root folder
			return nil

		case !strings.HasSuffix(info.Name(), ".json"):
			// Exclude not JSON files
			return nil
		}

		// Remaining: ./<contract_id>/<st_id>.json

		st := &sharetoken.T{}
		ps := strings.Split(path, pathSeparator)
		p_path := ps[initPathDepth:]
		err = t.m.Get(st, p_path...)

		if err != nil {
			// ToDo: log error
			// Halt only if file can't be moved
			return t.m.Rename(p_path, MalformedPath(p_path...))
		}

		return t.Load(st, p_path...)
	})

	return
}

func (t *T) add(st *sharetoken.T) (ps []string, err error) {
	k1, k2, k3 := t.keyf(st)

	if t.sts[k1] == nil {
		t.sts[k1] = st2map{}
	}

	if t.sts[k1][k2] == nil {
		t.sts[k1][k2] = st1map{}
	}

	if t.sts[k1][k2][k3] == nil {
		t.sts[k1][k2][k3] = st
	} else {
		err = DuplicateSTError
		return
	}

	ps = []string{k1, k3 + ".json"}

	return
}

// Load adds a sharetoken (st) to the map of accumulated sharetokens under the
// keys generated by t.keyf. It returns DuplicateSTError if this sharetoken
// was already seen.
func (t *T) Load(st *sharetoken.T, ps ...string) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	psX, err := t.add(st)

	if err != nil {
		// failed op, return
	} else if fsdir.PathEq(ps, psX) {
		// path checked, it's equal
	} else {
		// amending path
		err = t.m.Rename(ps, psX)
	}

	return
}

// Add adds a sharetoken (st) to the map of accumulated sharetokens under the
// keys generated by t.keyf. It returns DuplicateSTError if this sharetoken
// was already seen. Additionally, stores the correspondent file.
func (t *T) Add(st *sharetoken.T) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	path, err := t.add(st)

	if err == nil {
		err = t.m.Set(st, path...)
	}

	return
}

func (t *T) delete(st *sharetoken.T, handleFunc func(string, string) error) (err error) {
	k1, k2, k3 := t.keyf(st)

	switch {
	case t.sts[k1] == nil, t.sts[k1][k2] == nil, t.sts[k1][k2][k3] == nil:
		return
	}

	delete(t.sts[k1][k2], k3)
	err = handleFunc(k1, k3)

	if err != nil {
		return
	}

	if len(t.sts[k1][k2]) == 0 {
		delete(t.sts[k1], k2)
	}

	if len(t.sts[k1]) == 0 {
		delete(t.sts, k1)
	}

	return
}

// Del deletes a sharetoken (st) from the map of accumulated sharetokens
// under the keys generated by t.keyf. It can return errors from attempting to
// delete the file associated with the sharetoken on disk.
func (t *T) Del(st *sharetoken.T) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.delete(st, func(k1, k3 string) error {
		return t.m.Del(k1, k3+".json")
	})
}

// Exp expires a sharetoken (st) deleting it from the map of accumulated
// sharetokens under the keys generated by t.keyf. It can return errors from
// attempting to move the file associated with the sharetoken on disk.
func (t *T) Exp(st *sharetoken.T) (err error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.delete(st, func(k1, k3 string) error {
		path := []string{k1, k3 + ".json"}
		return t.m.Rename(path, ExpiredPath(path...))
	})
}

// Filter returns a list of sharetokens matching the given keys k1 and k2. An
// empty string for either of the keys is assumed to mean "for all values of
// this key".
func (t *T) Filter(k1, k2 string) (r []*sharetoken.T) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	switch {
	case k1 == "" && k2 == "":
		// for all k1, k2
		for _, m1 := range t.sts {
			for _, m2 := range m1 {
				for _, st := range m2 {
					r = append(r, st)
				}
			}
		}
	case k1 == "":
		// for all k1, some k2
		for _, m1 := range t.sts {
			for _, st := range m1[k2] {
				r = append(r, st)
			}
		}
	case k2 == "":
		// for some k1, all k2
		for _, m2 := range t.sts[k1] {
			for _, st := range m2 {
				r = append(r, st)
			}
		}
	default:
		// for some k1, some k2
		for _, st := range t.sts[k1][k2] {
			r = append(r, st)
		}
	}

	return
}

// SettlingAt returns a map with counts of sharetokens currently still being
// settled indexed by relay public key.
func (t *T) SettlingAt(rpk string, utime int64) map[string]int {
	r := map[string]int{}

	for _, st := range t.Filter("", rpk) {
		if st.IsSettlingAt(utime) {
			r[rpk]++
		}
	}

	return r
}
