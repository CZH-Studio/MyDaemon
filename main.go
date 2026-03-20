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

const (
	maxTailLines = 20
	logFileName  = "daemon.log"
)

type TailBuffer struct {
	lines []string
	mu    sync.Mutex
	max   int
}

func NewTailBuffer(n int) *TailBuffer {
	return &TailBuffer{
		lines: make([]string, 0, n),
		max:   n,
	}
}

func (t *TailBuffer) Add(line string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.lines) >= t.max {
		t.lines = t.lines[1:]
	}
	t.lines = append(t.lines, line)
}

func (t *TailBuffer) Get() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]string, len(t.lines))
	copy(out, t.lines)
	return out
}

func handleStream(reader io.Reader, writer io.Writer, tail *TailBuffer, wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(reader)

	// prevent too long lines
	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Fprintln(writer, line)
		tail.Add(line)
	}
}

func startProcess(cmd *exec.Cmd) (io.Reader, func(), error) {

	// PTY first
	if runtime.GOOS != "windows" {
		ptmx, err := pty.Start(cmd)
		if err == nil {
			cleanup := func() { _ = ptmx.Close() }
			return ptmx, cleanup, nil
		}

		// fallback
		fmt.Println("[daemon] PTY start failed, fallback to PIPE:", err)
	}

	// PIPE fallback
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	reader := io.MultiReader(stdout, stderr)

	cleanup := func() {}
	return reader, cleanup, nil
}

func parseArgs() []string {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mydaemon <command> [args...]")
		os.Exit(1)
		return nil
	} else {
		return os.Args[1:]
	}
}

func createLogFile() *os.File {
	logFile, err := os.Create(logFileName)
	if err != nil {
		fmt.Println("An error occurred when creating log file:", err)
		os.Exit(1)
	}
	return logFile
}

func main() {
	cmdArgs := parseArgs()
	logFile := createLogFile()
	defer logFile.Close()
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	tail := NewTailBuffer(maxTailLines)
	startTime := time.Now()
	reader, cleanup, err := startProcess(cmd)
	if err != nil {
		fmt.Println("An error occurred when starting process:", err)
		os.Exit(1)
	}
	defer cleanup()

	var wg sync.WaitGroup
	wg.Add(1)
	go handleStream(reader, io.MultiWriter(os.Stdout, logFile), tail, &wg)
	wg.Wait()
	err = cmd.Wait()
	duration := time.Since(startTime)

	fmt.Println("\n========== Process Result ==========")
	fmt.Println("Command:", cmdArgs)
	fmt.Println("Cost:", duration)

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	fmt.Println("Exit code:", exitCode)

	fmt.Println("\n===== Last 20 Lines of Output =====")
	for _, line := range tail.Get() {
		fmt.Println(line)
	}

	if runtime.GOOS == "windows" {
		fmt.Println("\n[Warning] Some programs in Windows may not flush stdout.")
		fmt.Println("It's recommended to use -u option in Python or flush stdout in the program itself.")
	}
}
