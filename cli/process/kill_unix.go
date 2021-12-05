// +build linux darwin

// Copyright (c) 2021 Wireleap
package process

import "golang.org/x/sys/unix"

// Term terminates the given PID gracefully.
func Term(pid int) error { return maybeSignal(pid, unix.SIGTERM) }

// Kill terminates the given PID forcefully.
func Kill(pid int) error { return maybeSignal(pid, unix.SIGKILL) }
