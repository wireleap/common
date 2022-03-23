//go:build linux || darwin
// +build linux darwin

// Copyright (c) 2022 Wireleap
package process

import "syscall"

const ReloadSignal = syscall.SIGUSR1

// Reload sends the reload signal (SIGUSR1) to the given PID.
func Reload(pid int) error { return maybeSignal(pid, ReloadSignal) }
