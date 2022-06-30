package client

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

func BroadcastQueueLocation(gid int64) {
	d := NewClient("User")
	services := d.GetServices()
	for _, ser := range services {
		server := ser
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(server.Key, "")
			c := client.NewXClient("User", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.GroupRequest{GroupId: gid}
			resp := &response.NilResponse{}
			_ = c.Call(context.Background(), "QueueLocation", req, resp)
		}()
	}
}
