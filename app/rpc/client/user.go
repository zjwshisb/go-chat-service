package client

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

func BroadcastQueueLocation(gid int64) {
	d := NewDiscovery("User")
	c := client.NewXClient("User", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()
	req := &request.GroupRequest{GroupId: gid}
	resp := &response.NilResponse{}
	_ = c.Broadcast(context.Background(), "QueueLocation", req, resp)
}
