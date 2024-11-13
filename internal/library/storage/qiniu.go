package storage

import (
	"context"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
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
func (s *qiniuAdapter) SaveUpload(ctx context.Context, file *ghttp.UploadFile, relativePath string) (string, error) {
	formUploader := storage.NewFormUploader(&storage.Config{})
	policy := storage.PutPolicy{
		Scope: s.bucket,
	}
	upToken := policy.UploadToken(qbox.NewMac(s.ak, s.sk))
	ret := storage.PutRet{}
	key := relativePath + "/" + random.RandString(32)
	f, err := file.Open()
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	err = formUploader.Put(ctx, &ret, upToken, key,
		f, file.Size, nil)
	if err != nil {
		return "", err
	}
	return key, nil

}
