package main

import "flag"

type Options struct {
	ConfigMode bool
	From       string
	Password   string
	To         string
	Command    []string
}

func ParseArgs() *Options {
	opts := &Options{}

	flag.BoolVar(&opts.ConfigMode, "config", false, "Save config")
	flag.StringVar(&opts.From, "from", "", "Sender email")
	flag.StringVar(&opts.Password, "password", "", "Email password")
	flag.StringVar(&opts.To, "to", "", "Receiver email")

	flag.Parse()
	opts.Command = flag.Args()

	return opts
}
