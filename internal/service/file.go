package service

import (
	"context"
	"gf-chat/api/v1"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IFile interface {
		trait.ICurd[model.CustomerChatFile]
		Insert(ctx context.Context, file *model.CustomerChatFile) (*model.CustomerChatFile, error)
		ToApi(file *model.CustomerChatFile) *v1.File
		FindAnd2Api(ctx context.Context, id any) (apiFile *v1.File, err error)
		Url(file *model.CustomerChatFile) string
		ThumbUrl(file *model.CustomerChatFile) string
		RemoveFile(ctx context.Context, file *model.CustomerChatFile) error
	}
)

var (
	localFile IFile
)

func File() IFile {
	if localFile == nil {
		panic("implement not found for interface IFile, forgot register?")
	}
	return localFile
}

func RegisterFile(i IFile) {
	localFile = i
}
