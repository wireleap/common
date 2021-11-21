// Copyright (c) 2021 Wireleap

package upgrade

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/blang/semver"
	"github.com/wireleap/common/cli/fsdir"
	"github.com/wireleap/common/cli/process"
)

func run(args ...string) (err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ExecutorArgs are the working parameters of the upgrade executor.
type ExecutorArgs struct {
	// Root directory of upgrade executor.
	Root fsdir.T
	// Source and destination binary paths relative to Root.
	SrcBin, DstBin string
	// Source and destination versions.
	SrcVer, DstVer semver.Version
}

type Executor func(ExecutorArgs) error

func ExecutorSimple(ea ExecutorArgs) (err error) {
	log.Printf("Running simple upgrade to %s...", ea.DstVer)
	log.Printf("stopping %s if running...", ea.SrcBin)
	var (
		pid     int
		pidfile = ea.SrcBin + ".pid"
		oldbin  = ea.DstBin + ".prev"
	)
	if err = ea.Root.Get(&pid, pidfile); err == nil && process.Exists(pid) {
		log.Printf("found old %s, pid %d", ea.DstBin, pid)
		if err = run(ea.SrcBin, "stop"); err != nil {
			err = fmt.Errorf("stopping old binary returned error: %s", err)
			return
		}
	}
	log.Printf("replacing %s with %s...", ea.DstBin, ea.SrcBin)
	if err = os.Rename(ea.DstBin, oldbin); err != nil {
		err = fmt.Errorf("renaming %s -> %s failed: %w", ea.DstBin, oldbin, err)
		return
	}
	if err = os.Rename(ea.SrcBin, ea.DstBin); err != nil {
		// attempt rollback
		err = fmt.Errorf("renaming %s -> %s failed: %s", ea.SrcBin, ea.DstBin, err)
		if e := os.Rename(oldbin, ea.DstBin); e != nil {
			err = fmt.Errorf("%w, additionally error when rolling back: %s", err, e)
		}
	}
	return
}

func ExecutorSupervised(ea ExecutorArgs) (err error) {
	// launch upgrade supervisor
	// if all goes well this binary will actually be terminated before this will return
	log.Printf("Running supervised upgrade to %s...", ea.DstVer)
	cmd := exec.Command(ea.SrcBin, "supervise-upgrade", ea.SrcVer.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("supervise-upgrade Start() returned error %w", err)
	}
	if err = cmd.Process.Release(); err != nil {
		return fmt.Errorf("supervise-upgrade Release() returned error %w", err)
	}
	return
}
