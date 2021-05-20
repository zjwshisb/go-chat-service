package file

import (
	"mime/multipart"
	"ws/configs"
)

type File struct {
	ThumbUrl string
	FullUrl string
	Path string
	Storage string
}

type Manager interface {
	Save(file *multipart.FileHeader, path string) (*File, error)
}

func Save(file *multipart.FileHeader, path string) (*File, error) {
	storage := configs.File.Storage
	var disk Manager
	switch storage {
	case "qiniu":
		disk = newQiniu()
	case "local":
		disk = NewLocal()
	default:
		disk = &Local{}
	}
	return disk.Save(file, path)
}