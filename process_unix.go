//go:build !windows

package main

import (
	"io"
	"os/exec"

	"github.com/creack/pty"
)

func StartProcess(cmd *exec.Cmd) (*ProcessHandle, error) {
	ptmx, err := pty.Start(cmd)
	if err == nil {
		return &ProcessHandle{
			Reader:  ptmx,
			Cleanup: func() { _ = ptmx.Close() },
			Wait: func() (int, error) {
				err := cmd.Wait()
				if err != nil {
					if e, ok := err.(*exec.ExitError); ok {
						return e.ExitCode(), nil
					}
					return -1, err
				}
				return 0, nil
			},
		}, nil
	}
	if cmd.Process != nil {
		return nil, err
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &ProcessHandle{
		Reader:  io.MultiReader(stdout, stderr),
		Cleanup: func() {},
		Wait: func() (int, error) {
			err := cmd.Wait()
			if err != nil {
				if e, ok := err.(*exec.ExitError); ok {
					return e.ExitCode(), nil
				}
				return -1, err
			}
			return 0, nil
		},
	}, nil
}
