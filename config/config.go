package config

import (
	"github.com/spf13/viper"
	"log"
)

func Setup(name string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(name)
	viper.AddConfigPath("./")
	viper.AddConfigPath("/")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("log file err: %v", err)
	}
}

func GetEnv() string {
	env := viper.GetString("App.Env")
	if env == "" {
		env = "local"
	}
	return env
}

func GetStoragePath() string {
	return GetWorkDir() + "/storage"
}

func GetWorkDir() string {
	workDir := viper.GetString("App.WorkDir")
	if workDir == "" {
		workDir = "./"
	}
	return workDir
}
