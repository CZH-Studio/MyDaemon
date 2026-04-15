package main

func main() {
	args := ParseArgs()
	switch {
	case args.Email:
		HandleEmail(args)
	case args.Config:
		HandleConfig(args)
	case args.Run:
		HandleRun(args)
	}
}
