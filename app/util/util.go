package util

import "github.com/spf13/viper"


func GetEnv() string {
	env :=  viper.GetString("App.Env")
	if env == "" {
		env = "local"
	}
	return env
}
func GetStoragePath() string  {
	return GetWorkDir() + "/storage"
}
func GetWorkDir() string {
	workDir := viper.GetString("App.WorkDir")
	if workDir == "" {
		workDir = "./"
	}
	return workDir
}