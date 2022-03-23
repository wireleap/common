// Copyright (c) 2022 Wireleap
package process

import (
	"errors"
	"log"
	"syscall"
)

var ReloadSignal = syscall.Signal(-1)

// Reload sends the reload signal (SIGUSR1) to the given PID. (NOT IMPLEMENTED YET)
func Reload(pid int) error {
	log.Printf("reload on windows not implemented yet, please use `wireleap restart` instead")
	return errors.New("reload on windows not implemented yet")
}
