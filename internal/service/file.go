package service

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/model"
	"gf-chat/internal/trait"
)

type (
	IFile interface {
		trait.ICurd[model.CustomerChatFile]
		Insert(ctx context.Context, file *model.CustomerChatFile) (*model.CustomerChatFile, error)
		ToApi(file *model.CustomerChatFile) *api.File
		FindAnd2Api(ctx context.Context, id any) (apiFile *api.File, err error)
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
