package config

import (
	"github.com/joho/godotenv"
	"log"
)

var (
	Config map[string]string
	DBConf db
)

func init()  {
	var err error
	Config, err = godotenv.Read()
	if err != nil {
		log.Fatal(nil)
	}
	DBConf = db{
		Host: Config["DB_HOST"],
		Port: Config["DB_PORT"],
		Database: Config["DB_DATABASE"],
		Username: Config["DB_USERNAME"],
		Password: Config["DB_PASSWORD"],
		Connection: Config["DB_CONNECTION"],
	}
	log.Println(DBConf)
}

type db struct {
	Host string
	Port string
	Database string
	Username string
	Password string
	Connection string
}
