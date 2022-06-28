package client

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

func SendMessage(id int64, server string) {
	d, _ := client.NewPeer2PeerDiscovery(server, "")
	c := client.NewXClient("Message", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()
	req := &request.SendMessageRequest{Id: id}
	resp := &response.OnlineResponse{}
	c.Call(context.Background(), "Send", req, resp)
}
