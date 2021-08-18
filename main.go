package main

import (
	"log"
	"os"
	"ws/app"
)

func main() {
	args := os.Args
	var command string
	if len(args) < 2 {
		command = "start"
	} else {
		command = args[1]
		if command[0:1] == "-" {
			command = "start"
		}
	}
	if command != "start" && command != "stop" && command != "restart" {
		log.Fatal("use ws start|stop -c=config.ini")
	}
	switch command {
	case "start":
		app.Start()
	case "stop":
		app.Stop()
	}
}
