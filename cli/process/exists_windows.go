// Copyright (c) 2021 Wireleap

package process

import "os"

// Exists checks if a given PID exists.
func Exists(pid int) bool {
	_, err := os.FindProcess(pid)
	return err == nil
}
