package chat

import (
	"context"
	rpcApi "gf-chat/api/chat/v1"
	v1 "gf-chat/api/chat/v1"
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	rpcApi.UnimplementedChatServer
}

func Register(s *grpcx.GrpcServer) {
	rpcApi.RegisterChatServer(s.Server, &Controller{})
}

func (*Controller) GetConnInfo(ctx context.Context, req *rpcApi.GetConnInfoRequest) (res *rpcApi.GetConnInfoReply, err error) {
	exist, platform := service.Chat().GetConnInfo(ctx, uint(req.CustomerId), uint(req.UserId), req.Type, true)
	return &rpcApi.GetConnInfoReply{
		Exist:    exist,
		Platform: platform,
	}, nil
}

func (*Controller) SendUserMessage(ctx context.Context, req *v1.SendUserMessageRequest) (res *v1.SendUserMessageReply, err error) {
	err = service.Chat().DeliveryUserMessage(ctx, uint(req.MsgId))
	if err != nil {
		return
	}
	return &v1.SendUserMessageReply{}, nil
}

func (*Controller) SendAdminMessage(ctx context.Context, req *v1.SendAdminMessageRequest) (res *v1.SendAdminMessageReply, err error) {
	err = service.Chat().DeliveryAdminMessage(ctx, uint(req.MsgId))
	if err != nil {
		return
	}
	return &v1.SendAdminMessageReply{}, nil
}

func (*Controller) NoticeRead(ctx context.Context, req *v1.NoticeReadRequest) (res *v1.NoticeReadReply, err error) {
	err = service.Chat().NoticeRead(ctx, uint(req.CustomerId), uint(req.UserId), gconv.Uints(req.MsgId), req.Type, true)
	if err != nil {
		return nil, err
	}
	return &v1.NoticeReadReply{}, nil
}

func (*Controller) GetOnlineUserIds(ctx context.Context, req *v1.GetOnlineUserIdsRequest) (res *v1.GetOnlineUserIdsReply, err error) {
	ids, err := service.Chat().GetOnlineUserIds(ctx, uint(req.CustomerId), req.Type, true)
	return &v1.GetOnlineUserIdsReply{Uid: gconv.Uint32s(ids)}, nil
}
