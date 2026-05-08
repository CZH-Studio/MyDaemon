//go:build windows

package main

import (
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/UserExistsError/conpty"
)

func buildCommandLine(cmd *exec.Cmd) string {
	var b strings.Builder
	for i, arg := range cmd.Args {
		if i > 0 {
			b.WriteByte(' ')
		}
		if strings.ContainsAny(arg, " \t\"") || arg == "" {
			b.WriteByte('"')
			b.WriteString(strings.ReplaceAll(arg, `"`, `\"`))
			b.WriteByte('"')
		} else {
			b.WriteString(arg)
		}
	}
	return b.String()
}

func StartProcess(cmd *exec.Cmd) (*ProcessHandle, error) {
	if conpty.IsConPtyAvailable() {
		cmdLine := buildCommandLine(cmd)
		var opts []conpty.ConPtyOption
		if cmd.Dir != "" {
			opts = append(opts, conpty.ConPtyWorkDir(cmd.Dir))
		}
		if len(cmd.Env) > 0 {
			opts = append(opts, conpty.ConPtyEnv(cmd.Env))
		}
		cpty, err := conpty.Start(cmdLine, opts...)
		if err == nil {
			return &ProcessHandle{
				Reader:  cpty,
				Cleanup: func() { cpty.Close() },
				Wait: func() (int, error) {
					code, err := cpty.Wait(context.Background())
					return int(code), err
				},
			}, nil
		}
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
