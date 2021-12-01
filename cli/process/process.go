// Copyright (c) 2021 Wireleap

package process

import (
	"os"
	"syscall"
)

func maybeSignal(pid int, sig os.Signal) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	// unix can only kill 0 if the process is actually alive
	return p.Signal(sig)
}

// NOTE: there is no distinction between Term and Kill on windows.

// Term terminates the given PID gracefully.
func Term(pid int) error { return maybeSignal(pid, syscall.SIGTERM) }

// Kill terminates the given PID forcefully.
func Kill(pid int) error { return maybeSignal(pid, os.Kill) }
