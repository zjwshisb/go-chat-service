package file

import (
	"context"
	"gf-chat/internal/model"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuAdapter struct {
	ak     string
	sk     string
	bucket string
	url    string
}

var (
	qiniu *QiniuAdapter
)

func init() {
	qiniu = newQiniu()
}

func newQiniu() *QiniuAdapter {
	ctx := gctx.New()
	ak, _ := g.Cfg().Get(ctx, "storage.qiniu.ak")
	sk, _ := g.Cfg().Get(ctx, "storage.qiniu.sk")
	bucket, _ := g.Cfg().Get(ctx, "storage.qiniu.bucket")
	baseUrl, _ := g.Cfg().Get(ctx, "storage.qiniu.url")
	return &QiniuAdapter{
		ak:     ak.String(),
		sk:     sk.String(),
		bucket: bucket.String(),
		url:    baseUrl.String(),
	}
}

func (s *QiniuAdapter) Url(path string) string {
	if len(path) >= 1 {
		first := path[0:1]
		if first == "/" {
			return s.url + path
		}
		return s.url + "/" + path
	}
	return ""
}
func (s *QiniuAdapter) Save(ctx context.Context, file *ghttp.UploadFile, relativePath string) (*model.SaveFileOutput, error) {
	cfg := &storage.Config{}
	policy := storage.PutPolicy{
		Scope: s.bucket,
	}
	formUploader := storage.NewFormUploader(cfg)
	mac := qbox.NewMac(s.ak, s.sk)
	upToken := policy.UploadToken(mac)
	ret := storage.PutRet{}
	key := relativePath + "/" + random.RandString(32)
	f, err := file.Open()
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return nil, err
	}
	err = formUploader.Put(ctx, &ret, upToken, key,
		f, file.Size, nil)
	if err != nil {
		return nil, err
	}
	return &model.SaveFileOutput{
		Url:  s.Url(key),
		Path: key,
	}, nil
}
