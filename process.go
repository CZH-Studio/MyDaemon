package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/creack/pty"
)

const logFileName = "daemon.log"

func CreateLogFile() *os.File {
	f, err := os.Create(logFileName)
	if err != nil {
		fmt.Println("Create log failed:", err)
		os.Exit(1)
	}
	return f
}

func Run(cmdArgs []string, logFile *os.File) (time.Duration, int, *TailBuffer) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	tail := NewTailBuffer(20)

	start := time.Now()

	reader, cleanup, err := startProcess(cmd)
	if err != nil {
		fmt.Println("Start process failed:", err)
		os.Exit(1)
	}
	defer cleanup()

	var wg sync.WaitGroup
	wg.Add(1)

	go handleStream(reader, io.MultiWriter(os.Stdout, logFile), tail, &wg)

	wg.Wait()

	err = cmd.Wait()
	duration := time.Since(start)

	exitCode := 0
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			exitCode = e.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return duration, exitCode, tail
}

func startProcess(cmd *exec.Cmd) (io.Reader, func(), error) {
	if runtime.GOOS != "windows" {
		ptmx, err := pty.Start(cmd)
		if err == nil {
			return ptmx, func() { _ = ptmx.Close() }, nil
		}
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	return io.MultiReader(stdout, stderr), func() {}, nil
}

func handleStream(reader io.Reader, writer io.Writer, tail *TailBuffer, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(writer, line)
		tail.Add(line)
	}
}

func FormatResult(cmdArgs []string, duration time.Duration, exitCode int, tail *TailBuffer) string {
	var b string

	b += "\n========== Process Result ==========\n"
	b += fmt.Sprintf("Command: %v\n", cmdArgs)
	b += fmt.Sprintf("Cost: %v\n", duration)
	b += fmt.Sprintf("Exit code: %d\n", exitCode)

	b += "\n===== Last 20 Lines =====\n"
	for _, line := range tail.Get() {
		b += line + "\n"
	}

	return b
}
