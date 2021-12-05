// Copyright (c) 2021 Wireleap

package process

import (
	"os"
)

func maybeSignal(pid int, sig os.Signal) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	// unix can only kill 0 if the process is actually alive
	return p.Signal(sig)
}
