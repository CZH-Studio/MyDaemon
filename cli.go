package main

import (
	"flag"
	"fmt"
	"os"
)

type Args struct {
	Email         bool
	EmailAdd      bool
	EmailAddFrom  string
	EmailAddPwd   string
	EmailAddTo    string
	EmailLs       bool
	EmailRm       bool
	Config        bool
	ConfigBuffer  int
	ConfigLog     string
	ConfigFrom    string
	ConfigSubject string
	Run           bool
	Command       []string
}

func ParseArgsEmailAdd(parser *Args, args []string) {
	fs := flag.NewFlagSet("email add", flag.ExitOnError)
	from := fs.String(Strings.OptEmailAddFrom, "", "Sender email")
	password := fs.String(Strings.OptEmailAddPwd, "", "SMTP password")
	to := fs.String(Strings.OptEmailAddTo, "", "Receiver email")
	fs.Parse(args)
	if *from == "" || *password == "" || *to == "" {
		fmt.Println(ErrorMissingFlags)
		fs.Usage()
		os.Exit(1)
	} else {
		parser.EmailAdd = true
		parser.EmailAddFrom = *from
		parser.EmailAddPwd = *password
		parser.EmailAddTo = *to
	}
}

func ParseArgsEmail(parser *Args, args []string) {
	if len(args) < 1 {
		fmt.Println(ErrorEmail)
		os.Exit(1)
	}
	parser.Email = true
	switch args[0] {
	case Strings.OptEmailAdd:
		ParseArgsEmailAdd(parser, args[1:])
	case Strings.OptEmailList:
		parser.EmailLs = true
	case Strings.OptEmailRemove:
		parser.EmailRm = true
	default:
		fmt.Println(ErrorEmail)
		os.Exit(1)
	}
}

func ParseArgsConfig(parser *Args, args []string) {
	fs := flag.NewFlagSet("config", flag.ExitOnError)
	buffer := fs.Int(Strings.OptConfigBuffer, 0, "Buffer size")
	log := fs.String(Strings.OptConfigLog, "", "Log file path")
	from := fs.String(Strings.OptConfigFrom, "", "Sender email")
	subject := fs.String(Strings.OptConfigSubject, "", "Email subject")
	fs.Parse(args)
	if *log == "" && *from == "" && *subject == "" && *buffer == 0 {
		fmt.Println(ErrorMissingFlags)
		fs.Usage()
		os.Exit(1)
	} else {
		parser.Config = true
		if *log != "" {
			parser.ConfigLog = *log
		}
		if *from != "" {
			parser.ConfigFrom = *from
		}
		if *subject != "" {
			parser.ConfigSubject = *subject
		}
		if *buffer > 0 && *buffer <= MAX_BUFFER {
			parser.ConfigBuffer = *buffer
		}
	}
}

func ParseArgsRun(parser *Args, args []string) {
	if len(args) < 1 {
		fmt.Println(ErrorMissingExecutable)
		os.Exit(1)
	}
	parser.Run = true
	parser.Command = args
}

func ParseArgs() *Args {
	usageString := UsageString()
	if len(os.Args) < 2 {
		fmt.Println(usageString)
		os.Exit(1)
	}
	parser := Args{}
	switch os.Args[1] {
	case Strings.OptEmail:
		ParseArgsEmail(&parser, os.Args[2:])
	case Strings.OptConfig:
		ParseArgsConfig(&parser, os.Args[2:])
	case Strings.OptRun:
		ParseArgsRun(&parser, os.Args[2:])
	default:
		fmt.Println(usageString)
		os.Exit(1)
	}
	return &parser
}
