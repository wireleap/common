// Copyright (c) 2021 Wireleap
package process

import "os"

// NOTE: there is no distinction between Term and Kill on windows.

// Term terminates the given PID gracefully.
func Term(pid int) error { return maybeSignal(pid, os.Kill) }

// Kill terminates the given PID forcefully.
func Kill(pid int) error { return maybeSignal(pid, os.Kill) }
