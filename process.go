package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

type ProcessHandle struct {
	Reader  io.Reader
	Cleanup func()
	Wait    func() (int, error)
}

func CreateLogFile(logFileName string) *os.File {
	f, err := os.Create(logFileName)
	if err != nil {
		fmt.Println(ErrorWritingFile)
		os.Exit(1)
	}
	return f
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
	ph, err := StartProcess(cmd)
	if err != nil {
		fmt.Println(ErrorLaunchingProcess, err)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go HandleStream(ph.Reader, io.MultiWriter(os.Stdout, logFile), tail, &wg)
	exitCode, err := ph.Wait()
	ph.Cleanup()
	wg.Wait()
	duration := time.Since(start)
	if err != nil {
		exitCode = -1
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
