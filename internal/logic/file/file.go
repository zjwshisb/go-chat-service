package file

import (
	"context"
	api "gf-chat/api/v1/backend"
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

func (s *sFile) SaveAndFill(ctx context.Context, file *model.CustomerChatFile) error {
	id, err := s.Save(ctx, file)
	if err != nil {
		return err
	}
	file.Id = uint(id)
	return nil
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
	return &api.File{
		Id:       file.Id,
		Path:     file.Path,
		Url:      storage.Disk(file.Disk).Url(file.Path),
		ThumbUrl: storage.Disk(file.Disk).ThumbUrl(file.Path),
		Name:     file.Name,
		Type:     "",
	}
}
