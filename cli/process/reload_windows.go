// Copyright (c) 2021 Wireleap
package process

import (
	"errors"
	"syscall"
)

var ReloadSignal = syscall.Signal(-1)

// Reload sends the reload signal (SIGUSR1) to the given PID. (NOT IMPLEMENTED YET)
func Reload(pid int) error { return errors.New("reload on windows not implemented yet") }
