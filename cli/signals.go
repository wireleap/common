// Copyright (c) 2021 Wireleap

package cli

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type SignalMap map[os.Signal]func() bool

func SignalLoop(hs SignalMap) {
	sigchan := make(chan os.Signal)

	for sig, _ := range hs {
		signal.Notify(sigchan, sig)
	}

	for sig := range sigchan {
		log.Printf("handling signal %d: %s...", sig, sig)

		if hs[sig]() {
			// unix convention -- exit code is 128 + terminating signal
			os.Exit(128 + int(sig.(syscall.Signal)))
		}
	}
}
