// Copyright (c) 2022 Wireleap

package process

import "os"

// Exists checks if a given PID exists.
func Exists(pid int) bool {
	p, err := os.FindProcess(pid)
	if err == nil {
		// release process handle
		// https://github.com/golang/go/issues/33814
		p.Release()
		return true
	} else {
		return false
	}
}
