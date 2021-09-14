package command

import (
	"flag"
	"log"
)

var (
	Command string
	ConfigFile string
	WithUser bool
)

func init()  {
	flag.StringVar(&Command, "m", "start" , "mode start|stop|migrate|fake|seeder")
	flag.StringVar(&ConfigFile,"c", "config.ini", "config file")
	flag.BoolVar(&WithUser, "u", false, "create user table when migrate")
	flag.Parse()
	if Command != "start" &&
		Command != "stop" &&
		Command != "fake" &&
		Command != "migrate" &&
		Command != "seeder" {
		log.Fatal("use ws start|stop|migrate|fake|seeder -c=config.ini -u")
	}
}
