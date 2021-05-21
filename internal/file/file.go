package file

import (
	"mime/multipart"
	"ws/configs"
)
type File struct {
	FullUrl string
	Path string
	Storage string
}

type Manager interface {
	Save(file *multipart.FileHeader, path string) (*File, error)
	Url(path string) string
}
func Disk(name string) Manager {
	switch name {
	case "qiniu":
		return newQiniu()
	case "local":
		return NewLocal()
	default:
		return NewLocal()
	}
}
func Save(file *multipart.FileHeader, path string) (*File, error) {
	storage := configs.File.Storage
	disk := Disk(storage)
	return disk.Save(file, path)
}