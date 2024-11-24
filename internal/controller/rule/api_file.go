package rule

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/library/storage"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
	"strings"
)

func init() {
	gvalid.RegisterRule("api-file", apiFile)
	gvalid.RegisterRule("api-file-if", apiFileIf)
}

// v:"file:image,5" or v:"file:image" or v:"file"
func file(ctx context.Context, in gvalid.RuleFuncInput) error {
	req := g.RequestFromCtx(ctx)
	uploadFile := req.GetUploadFile(in.Field)
	if uploadFile == nil {
		return nil
	}
	allowFileType := ""
	allowSize := 0
	params := getRuleParams(in.Rule)
	length := len(params)
	if length >= 1 {
		allowFileType = params[0]
	}
	if length >= 2 {
		allowSize = gconv.Int(params[1])
	}
	f, err := uploadFile.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	uploadFileType, err := storage.FileType(f)
	if err != nil {
		return err
	}
	defaultSize, exist := storage.DefaultFileSize[uploadFileType]
	if !exist {
		return gerror.New("unsupported file type")
	}
	if allowSize <= 0 {
		allowSize = defaultSize
	}
	if allowFileType != "" && strings.ToLower(uploadFileType) != strings.ToLower(allowFileType) {
		return gerror.New("unsupported file type")
	}
	if allowSize > 0 {
		if uploadFile.Size > int64(1024*1024*allowSize) {
			return gerror.New("file size too large")
		}
	}
	return nil
}

// v:"api-file-if:field,value"
func apiFileIf(ctx context.Context, in gvalid.RuleFuncInput) error {
	params := getRuleParams(in.Rule)
	if len(params) != 2 {
		return gerror.New("unsupported used for file-if rule")
	}
	field := params[0]
	value := params[1]
	data := in.Data.MapStrVar()
	if existV, ok := data[field]; ok {
		if existV.String() == value {
			return apiFile(ctx, in)
		} else {
			return nil
		}
	} else {
		return gerror.Newf("%s 必须是个有效的文件", in.Field)
	}
}

// v:"api-file"
func apiFile(ctx context.Context, in gvalid.RuleFuncInput) error {
	if in.Value.IsNil() {
		return nil
	}
	if in.Value.String() == "" {
		return nil
	}
	var apiFile *api.File
	err := in.Value.Scan(&apiFile)
	if err != nil {
		return gerror.Newf("%s 必须是个有效的文件", in.Field)
	}
	uid := service.AdminCtx().GetId(ctx)

	exists, err := service.File().Exists(ctx, do.CustomerChatFiles{
		Id:     apiFile.Id,
		FromId: uid,
	})
	if err != nil {
		return gerror.Newf("%s 必须是个有效的文件", in.Field)
	}
	if !exists {
		return gerror.Newf("%s 必须是个有效的文件", in.Field)
	}
	return nil
}
