package rule

import (
	"context"
	api "gf-chat/api/v1/backend"
	"gf-chat/internal/consts"
	"gf-chat/internal/library/storage"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
	"strings"
)

func init() {
	gvalid.RegisterRule("api-file", apiFile)
	gvalid.RegisterRule("file", file)
	gvalid.RegisterRule("file-dir", fileDir)
}

var (
	UnSupportFileError = gerror.New("不支持的文件类型")
)

// v:"file-dir"
func fileDir(ctx context.Context, in gvalid.RuleFuncInput) error {
	if in.Value.IsEmpty() {
		return nil
	}
	if in.Value.Uint() == 0 {
		return nil
	}
	parent, err := service.File().First(ctx, do.CustomerChatFiles{
		Id:         in.Value.Uint(),
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
		Type:       consts.FileTypeDir,
	})
	if err != nil {
		return err
	}
	if parent == nil {
		return gerror.NewCode(gcode.CodeValidationFailed, "上级目录不存在")
	}
	request := g.RequestFromCtx(ctx)
	request.SetCtxVar("file-dir", parent)
	return nil
}

// v:"file:image,5" or v:"file:image" or v:"file"
func file(ctx context.Context, in gvalid.RuleFuncInput) error {
	req := g.RequestFromCtx(ctx)
	uploadFile := req.GetUploadFile(strings.ToLower(in.Field))
	if uploadFile == nil {
		return nil
	}
	allowFileType := ""
	allowSize := 0
	params := parseRuleParams(in.Rule)
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
		return UnSupportFileError
	}
	if allowSize <= 0 {
		allowSize = defaultSize
	}
	if allowFileType != "" && strings.ToLower(uploadFileType) != strings.ToLower(allowFileType) {
		return UnSupportFileError
	}
	if allowSize > 0 {
		if uploadFile.Size > int64(1024*1024*allowSize) {
			return gerror.Newf("最多允许上传%dM的文件", allowSize)
		}
	}
	return nil
}

// v:"api-file" | v:"api-file:image" | v:"api-file:image,video"
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
		return UnSupportFileError
	}
	fileModel, err := service.File().First(ctx, do.CustomerChatFiles{
		Id:         apiFile.Id,
		CustomerId: service.AdminCtx().GetCustomerId(ctx),
	})
	if err != nil {
		return UnSupportFileError
	}
	params := parseRuleParams(in.Rule)
	if len(params) != 0 {
		for _, allowType := range params {
			if allowType == fileModel.Type {
				return nil
			}
		}
		return UnSupportFileError
	}
	return nil
}
