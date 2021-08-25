package file

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"ws/app/util"
	"ws/configs"
)

type qiniu struct {
	ak string
	sk string
	bucket string
	BaseUrl string
}

func newQiniu() *qiniu {
	return &qiniu{
		ak: configs.File.QiniuAk,
		sk: configs.File.QiniuSK,
		bucket: configs.File.QiniuBucket,
		BaseUrl: configs.File.QiniuUrl,
	}
}
func (qiniu *qiniu) Url(path string) string {
	first := path[0:1]
	if first == "/" {
		return qiniu.BaseUrl + path
	}
	return qiniu.BaseUrl + "/" + path
}
func (qiniu *qiniu) Save(file *multipart.FileHeader, relativePath string) (*File, error) {
	cfg := &storage.Config{}
	policy := storage.PutPolicy{
		Scope: qiniu.bucket,
	}
	formUploader := storage.NewFormUploader(cfg)
	mac := qbox.NewMac(qiniu.ak, qiniu.sk)
	upToken := policy.UploadToken(mac)
	ret := storage.PutRet{}
	key :=  relativePath + "/" + util.RandomStr(32)
	f, err := file.Open()
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	err = formUploader.Put(context.Background(), &ret, upToken, key,
		f, file.Size, nil)
	if err != nil {
		return nil, err
	}
	return &File{
		FullUrl: qiniu.Url(key),
		Path: key,
		Storage: configs.File.Storage,
	}, nil
}
