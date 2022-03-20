// Copyright (c) 2022 Wireleap

// Package fsdir provides an abstract interface to a directory on disk.
package fsdir

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// T is the type of a fsdir.
type T string

// mkdir creates a directory with the required parameters.
func mkdir(dir string) error {
	fi, err := os.Stat(dir)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)

			if err != nil {
				return fmt.Errorf("error while trying to mkdir %s: %w", dir, err)
			}

			fi, err = os.Stat(dir)

			if err != nil {
				return fmt.Errorf("error while trying to stat created %s: %w", dir, err)
			}
		} else {
			return fmt.Errorf("error while trying to stat %s: %w", dir, err)
		}
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s exists and is not a directory", dir)
	}

	return nil
}

// New creates a new fsdir.
func New(dir string) (T, error) {
	p, err := filepath.Abs(dir)

	if err != nil {
		return T(p), err
	}

	err = mkdir(p)

	if err != nil {
		return T(p), err
	}

	return T(p), nil
}

func (t T) Path(ps ...string) string {
	if len(ps) == 0 {
		return string(t)
	}

	return filepath.Join(append([]string{string(t)}, ps...)...)
}

// Get reads the file under the path ps and unmarshals its JSON contents into
// x.
func (t T) Get(x interface{}, ps ...string) error {
	p := t.Path(ps...)
	err := mkdir(filepath.Dir(p))

	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(p)

	if err != nil {
		return fmt.Errorf("fsdir.Get read: %w", err)
	}

	err = json.Unmarshal(b, x)

	if err != nil {
		return fmt.Errorf("error while trying to unmarshal %s: %w", p, err)
	}

	return nil
}

// Set marshals the x value into JSON and writes it to the the path ps.
func (t T) Set(x interface{}, ps ...string) error {
	p := t.Path(ps...)
	err := mkdir(filepath.Dir(p))

	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(x, "", "    ")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(p, b, 0644)

	if err != nil {
		return err
	}

	return nil
}

// Set marshals the x value into JSON and writes it to the the path ps.
func (t T) Rename(oldPS, newPS []string) (err error) {
	op := t.Path(oldPS...)
	np := t.Path(newPS...)
	err = mkdir(filepath.Dir(np))

	if err == nil {
		err = os.Rename(op, np)
	}

	return
}

// Del deletes the file or directory under a given path.
func (t T) Del(ps ...string) error {
	p := t.Path(ps...)
	return os.RemoveAll(p)
}

// Chmod changes the permissions of the file or directory under a given path.
func (t T) Chmod(mode os.FileMode, ps ...string) error {
	p := t.Path(ps...)
	return os.Chmod(p, mode)
}
