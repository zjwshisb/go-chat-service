package file

import (
	"mime/multipart"
	"os"
	"path"
	"ws/configs"
	"ws/util"
)

type Local struct {
	BaseUrl string
	StoragePath string
}

func NewLocal() *Local {
	return &Local{
		BaseUrl: configs.App.Url + configs.File.LocalPrefix,
		StoragePath: configs.File.LocalPath,
	}
}

func (local *Local) createDir(filePath string)  error  {
	if !local.isExist(filePath) {
		err := os.MkdirAll(filePath,os.ModePerm)
		return err
	}
	return nil
}

func (local *Local) isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (local *Local) Save(file *multipart.FileHeader, relativePath string) (*File, error) {
	ff, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = ff.Close()
	}()
	ext := path.Ext(file.Filename)
	fileByte := make([]byte, file.Size)
	_ , err = ff.Read(fileByte)
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
		Path: relativeName,
		ThumbUrl: local.BaseUrl + "/" + relativeName,
		FullUrl: local.BaseUrl + "/" + relativeName,
		Storage: "local",
	}, nil
}