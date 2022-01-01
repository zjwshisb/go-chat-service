package util

import "github.com/spf13/viper"

func GetEnv() string {
	return  viper.GetString("App.Env")
}
