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

func CreateLogFile(logFileName string) *os.File {
	f, err := os.Create(logFileName)
	if err != nil {
		fmt.Println(ErrorWritingFile)
		os.Exit(1)
	}
	return f
}

func StartProcess(cmd *exec.Cmd) (io.Reader, func(), error) {
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

func HandleStream(reader io.Reader, writer io.Writer, tail *TailBuffer, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(writer, line)
		tail.Add(line)
	}
}

func Run(cmdArgs []string, logFile *os.File, bufferSize int) (time.Duration, int, *TailBuffer) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	tail := NewTailBuffer(bufferSize)
	start := time.Now()
	reader, cleanup, err := StartProcess(cmd)
	if err != nil {
		fmt.Println(ErrorLaunchingProcess, err)
		os.Exit(1)
	}
	defer cleanup()
	var wg sync.WaitGroup
	wg.Add(1)
	go HandleStream(reader, io.MultiWriter(os.Stdout, logFile), tail, &wg)
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

func HandleRun(args *Args) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Println(ErrorReadingFile)
		os.Exit(1)
	}
	logFile := CreateLogFile(config.Value.LogName)
	defer logFile.Close()
	duration, exitCode, tail := Run(args.Command, logFile, config.Value.BufferSize)
	result := FormatResult(args.Command, duration, exitCode, tail, config.Value.BufferSize)
	// send email
	for _, email := range config.Emails {
		err = SendEmail(email, config.Value.FromTitle, config.Value.SubjectTitle, result)
		if err != nil {
			fmt.Println(ErrorEmailSending, err)
		} else {
			fmt.Printf(PromptEmailSendingSucceed, email.From, email.To)
		}
	}
}
