package main

import (
	"os"
	"os/exec"
	"ws/app"
	"ws/args"
)
func daemonize(args ...string)  {
	var arg []string
	if len(args) > 1 {
		arg = args[1:]
	}
	cmd := exec.Command(args[0], arg...)
	cmd.Env = os.Environ()
	cmd.Start()
}
func main() {
	if !args.Daemonized {
		app.Start()
	} else {
		args := os.Args
		daemonize(args...)
	}
}
