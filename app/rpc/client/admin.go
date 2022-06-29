package client

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

func BroadcastOnlineAdmin(groupId int64) {
	d := NewClient("Admin")
	services := d.GetServices()
	for _, ser := range services {
		server := ser
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(server.Key, "")
			c := client.NewXClient("Admin", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.GroupRequest{GroupId: groupId}
			resp := &response.NilResponse{}
			_ = c.Call(context.Background(), "OnlineAdmin", req, resp)
		}()
	}
}

func NoticeUserTransfer(id int64, server string) {
	d, _ := client.NewPeer2PeerDiscovery(server, "")
	c := client.NewXClient("Admin", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()
	req := &request.IdRequest{Id: id}
	resp := &response.NilResponse{}
	_ = c.Call(context.Background(), "UserTransfer", req, resp)
}

func NoticeUserOnline(uid int64, server string) {
	d, _ := client.NewPeer2PeerDiscovery(server, "")
	c := client.NewXClient("Admin", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()
	req := &request.IdRequest{Id: uid}
	resp := &response.NilResponse{}
	_ = c.Call(context.Background(), "UserOnline", req, resp)
}

func NoticeUserOffLine(uid int64, server string) {
	d, _ := client.NewPeer2PeerDiscovery(server, "")
	c := client.NewXClient("Admin", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()
	req := &request.IdRequest{Id: uid}
	resp := &response.NilResponse{}
	_ = c.Call(context.Background(), "UserOffline", req, resp)
}

func BroadcastWaitingUser(groupId int64) {
	d := NewClient("Admin")
	services := d.GetServices()
	for _, ser := range services {
		server := ser
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(server.Key, "")
			c := client.NewXClient("Admin", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.GroupRequest{GroupId: groupId}
			resp := &response.NilResponse{}
			_ = c.Call(context.Background(), "WaitingUser", req, resp)
		}()
	}
}
