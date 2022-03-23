//go:build linux || darwin
// +build linux darwin

// Copyright (c) 2022 Wireleap
package process

import (
	"os"
	"syscall"
)

// Exists checks if a given PID exists.
func Exists(pid int) bool {
	// unix can get a permission error if the process is actually alive
	err := maybeSignal(pid, syscall.Signal(0))
	return err == nil || os.IsPermission(err)
}
