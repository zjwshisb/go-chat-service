package storage

import (
	"context"
	"gf-chat/internal/consts"
	"gf-chat/internal/model"
	"gf-chat/internal/model/entity"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"strings"
)

type qiniuAdapter struct {
	ak     string
	sk     string
	bucket string
	url    string
}

func newQiniu() *qiniuAdapter {
	ctx := gctx.New()
	ak, _ := g.Cfg().Get(ctx, "storage.qiniu.ak")
	sk, _ := g.Cfg().Get(ctx, "storage.qiniu.sk")
	bucket, _ := g.Cfg().Get(ctx, "storage.qiniu.bucket")
	baseUrl, _ := g.Cfg().Get(ctx, "storage.qiniu.url")
	return &qiniuAdapter{
		ak:     ak.String(),
		sk:     sk.String(),
		bucket: bucket.String(),
		url:    baseUrl.String(),
	}
}

func (s *qiniuAdapter) Url(path string) string {
	if len(path) >= 1 {
		first := path[0:1]
		if first == "/" {
			return s.url + path
		}
		return s.url + "/" + path
	}
	return ""
}
func (s *qiniuAdapter) ThumbUrl(path string) string {
	return s.Url(path)
}
func (s *qiniuAdapter) Delete(path string) error {
	m := storage.NewBucketManager(&auth.Credentials{
		AccessKey: s.ak,
		SecretKey: []byte(s.sk),
	}, nil)
	return m.Delete(s.bucket, path)
}
func (s *qiniuAdapter) SaveUpload(ctx context.Context, file *ghttp.UploadFile, relativePath string) (*model.CustomerChatFile, error) {
	formUploader := storage.NewFormUploader(&storage.Config{})
	policy := storage.PutPolicy{
		Scope: s.bucket,
	}
	upToken := policy.UploadToken(qbox.NewMac(s.ak, s.sk))
	ret := storage.PutRet{}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	if !strings.HasSuffix(relativePath, pathSeparator) {
		relativePath = relativePath + pathSeparator
	}
	ext := gfile.Ext(file.Filename)
	path := relativePath + random.RandString(32) + "." + ext
	err = formUploader.Put(ctx, &ret, upToken, path,
		f, file.Size, nil)
	if err != nil {
		return nil, err
	}
	return &model.CustomerChatFile{
		CustomerChatFiles: entity.CustomerChatFiles{
			Path: path,
			Disk: consts.StorageQiniu,
			Name: file.Filename,
		},
	}, nil

}
