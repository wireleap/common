// Copyright (c) 2021 Wireleap

package cli

import (
	"os"
	"os/exec"
)

func RunChild(args ...string) (err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
