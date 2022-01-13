package config

import (
	"github.com/spf13/viper"
	"log"
)

func Setup(name string)  {

	viper.SetConfigType("yaml")
	viper.SetConfigName(name)
	viper.AddConfigPath("./")
	viper.AddConfigPath("/")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("log file err: %v", err)
	}

}
