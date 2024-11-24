package rule

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gvalid"
	"strconv"
)

func init() {
	name := "image"
	gvalid.RegisterRule(name, imageRule)
}

func imageRule(ctx context.Context, in gvalid.RuleFuncInput) error {
	name := in.Rule
	var maxM int64 = 5
	if len(name) >= 6 {
		paramsStr := name[6:]
		params := gstr.Explode(",", paramsStr)
		length := len(params)
		if length == 1 {
			maxStr := params[0]
			maxInt, err := strconv.Atoi(maxStr)
			if err == nil {
				maxM = int64(maxInt)
			}
		}
	}
	req := g.RequestFromCtx(ctx)
	file := req.GetUploadFile("file")
	// bug,没法使用下面的方法获取,所以文件的key必须是file
	//file, ok := in.Value.Val().(*ghttp.UploadFile)
	if file != nil {
		if file.Size > maxM*1024*1024 {
			return gerror.Newf("图片大小必须小于%dM", maxM)
		}
	}
	return nil
}
