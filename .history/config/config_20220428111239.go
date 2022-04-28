package config

import (
	"log"

	"github.com/spf13/viper"
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
func IsCluster() bool {
	return viper.GetBool("App.Cluster")
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
