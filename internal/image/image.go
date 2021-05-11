package image

import (
	"log"
	"os"
	"ws/configs"
)

var (
	BasePath = "./storage/images"
	AvatarDIR = "/avatar"
	ChatDir =  "/chat"
)
func Setup()  {
	dirs := [3]string{
		BasePath,
		BasePath + AvatarDIR,
		BasePath + ChatDir,
	}
	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir(dir, 0755)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
func Url(path string)  string {
	return configs.App.Url + "/images" + path
}
