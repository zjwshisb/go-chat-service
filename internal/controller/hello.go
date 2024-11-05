package controller

import (
	"context"
	"fmt"
	"gf-chat/internal/dao"
	"gf-chat/internal/model/entity"
	"github.com/gogf/gf/v2/frame/g"

	"gf-chat/api/v1"
)

var (
	Hello = cHello{}
)

type cHello struct{}

func (c *cHello) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	var message []*entity.CustomerChatAutoMessages
	err = dao.CustomerChatMessages.Ctx(ctx).Fields(dao.CustomerChatMessages.Columns()).
		Scan(&message)
	fmt.Println(message)
	g.DB("")
	g.RequestFromCtx(ctx).Response.Writeln("Hello World!")
	return
}
