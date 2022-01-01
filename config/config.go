package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func Setup(name string)  {

	viper.SetConfigType("yaml")
	viper.SetConfigName(name)
	viper.AddConfigPath("/")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("config file: %w \n", err))
	}
}
