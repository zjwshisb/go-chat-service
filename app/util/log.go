package util

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/golang-module/carbon"
	"github.com/spf13/viper"
)

func Debug(name string, i interface{})  {
	if viper.GetString("App.Env") == "local" {
		time := carbon.Now().ToDateTimeString()
		f := color.New(color.FgCyan, color.Bold)
		f.Printf("[Debug][%s]%s %+v \n",time, name, i)
	}
}
func DebugWebsocket(t string, i interface{})  {
	name := fmt.Sprintf("[webscoket][%s]", t)
	Debug(name, i)
}
