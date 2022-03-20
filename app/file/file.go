package file

import (
	"github.com/spf13/viper"
	"mime/multipart"
)

const (
	StorageQiniu = "qiniu"
	StorageLocal = "local"
)

var diskQiniu *qiniu
var diskLocal *local

func Setup() {
	diskLocal = NewLocal()
	diskQiniu = NewQiniu()
}

type File struct {
	FullUrl string
	Path    string
	Storage string
}

type Manager interface {
	Save(file *multipart.FileHeader, path string) (*File, error)
	Url(path string) string
}

func Disk(name string) Manager {
	switch name {
	case StorageQiniu:
		return diskQiniu
	case StorageLocal:
		return diskLocal
	default:
		return diskLocal
	}
}

func Save(file *multipart.FileHeader, path string) (*File, error) {
	def := viper.GetString("File.Storage")
	disk := Disk(def)
	return disk.Save(file, path)
}
