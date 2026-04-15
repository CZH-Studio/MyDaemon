package main

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"
)

var Strings = struct {
	ProgramName      string
	ProgramVersion   string
	OptEmail         string
	OptEmailAdd      string
	OptEmailAddFrom  string
	OptEmailAddPwd   string
	OptEmailAddTo    string
	OptEmailList     string
	OptEmailRemove   string
	OptConfig        string
	OptConfigBuffer  string
	OptConfigLog     string
	OptConfigFrom    string
	OptConfigSubject string
	OptRun           string
}{
	ProgramName:      "mydaemon",
	ProgramVersion:   "1.1",
	OptEmail:         "email",
	OptEmailAdd:      "add",
	OptEmailAddFrom:  "from",
	OptEmailAddPwd:   "pwd",
	OptEmailAddTo:    "to",
	OptEmailList:     "ls",
	OptEmailRemove:   "rm",
	OptConfig:        "config",
	OptConfigBuffer:  "buffer",
	OptConfigLog:     "log",
	OptConfigFrom:    "from",
	OptConfigSubject: "subject",
	OptRun:           "run",
}

const UsageTemplate = `Usage:
	{{.ProgramName}} {{.OptEmail}} {{.OptEmailAdd}}
		--{{.OptEmailAddFrom}} xxx@example.com
		--{{.OptEmailAddPwd}} SMTP password
		--{{.OptEmailAddTo}} yyy@example.com
	{{.ProgramName}} {{.OptEmail}} {{.OptEmailList}}
	{{.ProgramName}} {{.OptEmail}} {{.OptEmailRemove}}

	{{.ProgramName}} {{.OptConfig}}
		--{{.OptConfigBuffer}} 20
		--{{.OptConfigLog}} logFileName.log
		--{{.OptConfigFrom}} MyDaemon
		--{{.OptConfigSubject}} EmailSubject

	{{.ProgramName}} {{.OptRun}} <command> [args ...]
`
const ErrorEmail = "Email requires subcommand: add | ls | remove"
const ErrorMissingFlags = "Missing required flags!"
const ErrorMissingExecutable = "Missing executable command!"
const ErrorReadingFile = "Error occurred when reading file!"
const ErrorWritingFile = "Error occurred when writing file!"
const ErrorInvalidIndex = "Invalid index!"
const ErrorLaunchingProcess = "Error occurred when launching process: "
const ErrorEmailSending = "Error occurred when sending email: "
const PromptEmailSendingSucceed = "Email sent successfully! (%s -> %s)\n"
const PromptRemovingEmail = "Select index to execute removal: "
const MAX_BUFFER = 1024

func UsageString() string {
	tpl, err := template.New("usage").Parse(UsageTemplate)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, Strings)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func FormatResult(cmdArgs []string, duration time.Duration, exitCode int, tail *TailBuffer, bufferSize int) string {
	var b strings.Builder
	b.WriteString("\n========== Process Result ==========\n")
	fmt.Fprintf(&b, "Command: %s\n", strings.Join(cmdArgs, " "))
	fmt.Fprintf(&b, "Cost: %v\n", duration)
	fmt.Fprintf(&b, "Exit code: %d\n", exitCode)
	fmt.Fprintf(&b, "\n===== Last %d Lines =====\n", bufferSize)
	for _, line := range tail.Get() {
		b.WriteString(line + "\n")
	}
	return b.String()
}
