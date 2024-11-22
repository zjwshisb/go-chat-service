package rule

import (
	"context"
	"gf-chat/api"
	"gf-chat/internal/model/do"
	"gf-chat/internal/service"
	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
	"strings"
)

func init() {
	gvalid.RegisterRule("api-file", apiFile)
	gvalid.RegisterRule("api-file-if", apiFileIf)
}

var defaultFileSize = map[string]int{
	"image": 5 * 1024 * 1024,
	"video": 10 * 1024 * 1024,
	"audio": 5 * 1024 * 1024,
}

// v:"file:image,5" or v:"file:image" or v:"file"
func file(ctx context.Context, in gvalid.RuleFuncInput) error {
	req := g.RequestFromCtx(ctx)
	uploadFile := req.GetUploadFile(in.Field)
	if uploadFile == nil {
		return nil
	}
	name := in.Rule
	allowFileType := ""
	allowSize := 0
	if len(name) >= 5 {
		paramsStr := name[5:]
		params := gstr.Explode(",", paramsStr)
		length := len(params)
		if length >= 1 {
			allowFileType = params[0]
		}
		if length >= 2 {
			allowSize = gconv.Int(params[1])
		}
	}
	f, err := uploadFile.Open()
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	mimeType := fileutil.MiMeType(f)
	index := strings.Index(mimeType, "/")
	if index < 0 {
		return gerror.New("unsupported file")
	}
	uploadFileType := mimeType[:index]
	defaultSize, exist := defaultFileSize[uploadFileType]
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
	rule := in.Rule
	index := strings.Index(rule, ":")
	if index == -1 {
		return gerror.New("unsupported used for file-if rule")
	}
	keyValue := rule[index+1:]
	keyValueArr := strings.Split(keyValue, ",")
	if len(keyValueArr) != 2 {
		return gerror.New("unsupported used for file-if rule")
	}
	field := keyValueArr[0]
	value := keyValueArr[1]
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
