package file

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/dao"
	"gf-chat/internal/library/storage"
	"gf-chat/internal/model"
	"gf-chat/internal/service"
	"gf-chat/internal/trait"
)

func init() {
	service.RegisterFile(&sFile{
		Curd: trait.Curd[model.CustomerChatFile]{
			Dao: &dao.CustomerChatFiles,
		},
	})
}

type sFile struct {
	trait.Curd[model.CustomerChatFile]
}

func (s *sFile) Insert(ctx context.Context, file *model.CustomerChatFile) (*model.CustomerChatFile, error) {
	id, err := s.Save(ctx, file)
	if err != nil {
		return nil, err
	}
	file.Id = uint(id)
	return file, nil
}

func (s *sFile) Url(file *model.CustomerChatFile) string {
	return storage.Disk(file.Disk).Url(file.Path)
}

func (s *sFile) ThumbUrl(file *model.CustomerChatFile) string {
	return storage.Disk(file.Disk).ThumbUrl(file.Path)
}

func (s *sFile) FindAnd2Api(ctx context.Context, id any) (apiFile *api.File, err error) {
	file, err := s.Find(ctx, id)
	if err != nil {
		return
	}
	apiFile = s.ToApi(file)
	return
}
func (s *sFile) ToApi(file *model.CustomerChatFile) *api.File {
	apiFile := &api.File{
		Id:   file.Id,
		Path: file.Path,
		Name: file.Name,
		Type: file.Type,
	}
	if file.Type != consts.FileTypeDir {
		apiFile.Url = s.Url(file)
		apiFile.ThumbUrl = s.ThumbUrl(file)
	}
	return apiFile
}
