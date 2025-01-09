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

func (*Controller) SendMessage(ctx context.Context, req *v1.SendMessageRequest) (res *v1.NilReply, err error) {
	err = service.Chat().DeliveryMessage(ctx, uint(req.MsgId), req.Type, true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}

func (*Controller) UpdateAdminSetting(ctx context.Context, req *v1.UpdateAdminSettingRequest) (res *v1.NilReply, err error) {
	err = service.Chat().UpdateAdminSetting(ctx, uint(req.Id), true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}

func (*Controller) BroadcastWaitingUser(ctx context.Context, req *v1.BroadcastWaitingUserRequest) (res *v1.NilReply, err error) {
	err = service.Chat().BroadcastWaitingUser(ctx, uint(req.CustomerId), true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}

func (*Controller) NoticeTransfer(ctx context.Context, req *v1.NoticeTransferRequest) (res *v1.NilReply, err error) {
	err = service.Chat().NoticeTransfer(ctx, uint(req.CustomerId), uint(req.AdminId), true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}

func (*Controller) NoticeUserOnline(ctx context.Context, req *v1.NoticeUserOnlineRequest) (res *v1.NilReply, err error) {
	err = service.Chat().NoticeUserOnline(ctx, uint(req.UserId), req.Platform, true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}

func (*Controller) NoticeUserOffline(ctx context.Context, req *v1.NoticeUserOfflineRequest) (res *v1.NilReply, err error) {
	err = service.Chat().NoticeUserOffline(ctx, uint(req.UserId), true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}

func (*Controller) BroadcastOnlineAdmins(ctx context.Context, req *v1.BroadcastOnlineAdminsRequest) (res *v1.NilReply, err error) {
	err = service.Chat().BroadcastOnlineAdmins(ctx, uint(req.CustomerId), true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}

func (*Controller) BroadcastQueueLocation(ctx context.Context, req *v1.BroadcastQueueLocationRequest) (res *v1.NilReply, err error) {
	err = service.Chat().BroadcastQueueLocation(ctx, uint(req.CustomerId), true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}

func (*Controller) NoticeRepeatConnect(ctx context.Context, req *v1.NoticeRepeatConnectRequest) (res *v1.NilReply, err error) {
	err = service.Chat().NoticeRepeatConnect(ctx, uint(req.UserId), uint(req.CustomerId), req.Type, req.NewUid, true)
	if err != nil {
		return
	}
	return &v1.NilReply{}, nil
}
