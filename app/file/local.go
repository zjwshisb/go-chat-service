package file

import (
	"github.com/spf13/viper"
	"mime/multipart"
	"os"
	"path"
	"ws/app/util"
)

type local struct {
	BaseUrl     string
	StoragePath string
}

func NewLocal() *local {
	return &local{
		BaseUrl:     viper.GetString("App.Url") + viper.GetString("File.LocalPrefix"),
		StoragePath: viper.GetString("File.LocalPath"),
	}
}

func (local *local) Url(path string) string {
	first := path[0:1]
	if first == "/" {
		return local.BaseUrl + path
	}
	return local.BaseUrl + "/" + path
}

func (local *local) createDir(filePath string) error {
	if !local.isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

func (local *local) isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
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
	filename := util.RandomStr(30) + ext
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

	err = local.createDir(fullPath)
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
		Path:     relativeName,
		FullUrl:  local.Url(relativeName),
		Storage:  StorageLocal,
	}, nil
}
