package args

import "flag"

var ConfigFile string
var Daemonized bool

func init()  {
	flag.StringVar(&ConfigFile, "c", "config.ini", "config file")
	flag.BoolVar(&Daemonized, "d", false, "is daemonized")
}
