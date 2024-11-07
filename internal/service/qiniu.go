// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gf-chat/internal/model"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IQiniu interface {
		Form(path string) *model.ImageFiled
		Url(path string) string
		Save(ctx context.Context, file *ghttp.UploadFile, relativePath string) (*model.SaveFileOutput, error)
	}
)

var (
	localQiniu IQiniu
)

func Qiniu() IQiniu {
	if localQiniu == nil {
		panic("implement not found for interface IQiniu, forgot register?")
	}
	return localQiniu
}

func RegisterQiniu(i IQiniu) {
	localQiniu = i
}
