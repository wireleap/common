// +build linux darwin

// Copyright (c) 2021 Wireleap

package process

import (
	"os"

	"golang.org/x/sys/unix"
)

func Writable(p string) bool {
	err := unix.Access(p, unix.W_OK)
	return err == nil || os.IsNotExist(err)
}
