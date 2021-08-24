package main

import (
	"ws/app"
	"ws/command"
)


func main() {
	switch command.Command {
	case "start":
		app.Start()
	case "stop":
		app.Stop()
	case "migrate":
		app.Migrate()
	case "fake":
		app.Fake()
	case "seeder":
		app.Seeder()
	}
}
