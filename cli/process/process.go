// Copyright (c) 2021 Wireleap

package process

import (
	"os"
	"syscall"
)

func maybeSignal(pid int, sig os.Signal) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		// windows will fail to find non-existent pid
		return err
	}
	// unix can only kill 0 if the process is actually alive
	return p.Signal(sig)
}

// Exists checks if a given PID exists in a portable way.
func Exists(pid int) bool {
	// unix can get a permission error if the process is actually alive
	err := maybeSignal(pid, syscall.Signal(0))
	return err == nil || os.IsPermission(err)
}

// NOTE: there is no distinction between Term and Kill on windows.

// Term terminates the given PID gracefully.
func Term(pid int) error { return maybeSignal(pid, syscall.SIGTERM) }

// Kill terminates the given PID forcefully.
func Kill(pid int) error { return maybeSignal(pid, os.Kill) }
