package file

import (
	"github.com/duke-git/lancet/v2/random"
	"mime/multipart"
	"os"
	"path"
	"ws/config"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/spf13/viper"
)

type local struct {
	BaseUrl     string
	StoragePath string
}

const prefix = "assets"

func NewLocal() *local {
	storagePath := config.GetStoragePath() + "/" + prefix
	if !fileutil.IsDir(storagePath) {
		err := os.MkdirAll(storagePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return &local{
		BaseUrl:     viper.GetString("App.Url") + "/" + prefix,
		StoragePath: storagePath,
	}
}

func (local *local) Url(path string) string {
	first := path[0:1]
	if first == "/" {
		return local.BaseUrl + path
	}
	return local.BaseUrl + "/" + path
}

func (local *local) Save(file *multipart.FileHeader, relativePath string) (*File, error) {
	ff, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = ff.Close()
	}()
	ext := path.Ext(file.Filename)
	fileByte := make([]byte, file.Size)
	_, err = ff.Read(fileByte)
	if err != nil {
		return nil, err
	}
	filename := random.RandString(32) + ext
	var fullPath string
	var relativeName string
	if relativePath != "" {
		fullPath = local.StoragePath + "/" + relativePath
		relativeName = relativePath + "/" + filename
	} else {
		fullPath = local.StoragePath
		relativeName = filename
	}
	fullName := fullPath + "/" + filename

	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	saveFile, err := os.Create(fullName)
	defer func() {
		_ = saveFile.Close()
	}()
	_, err = saveFile.Write(fileByte)
	if err != nil {
		return nil, err
	}
	return &File{
		Path:    relativeName,
		FullUrl: local.Url(relativeName),
		Storage: StorageLocal,
	}, nil
}
