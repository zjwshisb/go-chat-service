package util

import (
	"log"
	"os"
	"ws/config"
)

func Asset(path string)  string {
	return config.Http.Url + "/assets" + path
}
func AssetPath(path string) string {
	return BasePath() + "/storage/assets"
}
func StoragePath() string {
	return BasePath() + "/storage"
}
func BasePath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return path
}