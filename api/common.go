package api

import (
	"gf-chat/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ListRes[T any] struct {
	NormalRes[[]T]
	Total int `json:"total"`
}

type NilRes = NormalRes[any]

type OptionRes struct {
	NormalRes[[]model.Option]
}

type NormalRes[T any] struct {
	Code    int  `json:"code"    dc:"Error code"`
	Data    T    `json:"data"    dc:"Result data for certain request according API definition"`
	Success bool `json:"success" dc:"Is Success"`
}

type FailRes struct {
	Code    int    `json:"code"    dc:"Error code"`
	Success bool   `json:"success" dc:"Is Success"`
	Message string `json:"message" dc:"错误消息"`
}

type ImageReq struct {
	g.Meta `path:"/images" tags:"后台图片上传" mine:"multipart/form-data" method:"post" summary:"上传图片"`
	Path   string            `json:"path" p:"path" v:"required" dc:"文件存储路径"`
	File   *ghttp.UploadFile `json:"file" p:"file" type:"file" v:"image" dc:"文件"`
}

type ImageRes struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

func NewOptionResp(options []model.Option) *OptionRes {
	return &OptionRes{
		NormalRes: NormalRes[[]model.Option]{
			Code:    0,
			Success: true,
			Data:    options,
		},
	}
}

func NewFailResp(message string, code int) *FailRes {
	return &FailRes{
		Code:    code,
		Success: false,
		Message: message,
	}
}

func NewListResp[T any](items []T, total int) *ListRes[T] {
	return &ListRes[T]{
		NormalRes: NormalRes[[]T]{
			Success: true,
			Data:    items,
			Code:    0,
		},
		Total: total,
	}
}
func NewNilResp() *NilRes {
	return &NilRes{
		Success: true,
		Data:    nil,
		Code:    0,
	}
}
func NewResp[T any](data T) *NormalRes[T] {
	return &NormalRes[T]{
		Success: true,
		Data:    data,
		Code:    0,
	}
}
