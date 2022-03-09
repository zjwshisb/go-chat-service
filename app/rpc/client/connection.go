package client

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"sync"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

func ClientTotal(groupId int64, types string) int64 {
	d := NewClient("Connection")
	services := d.GetServices()
	var total int64
	result := make(chan int64, len(services))
	var wg  sync.WaitGroup
	for _, service:= range services{
		wg.Add(1)
		ser := service
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(ser.Key, "")
			c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.GroupRequest{GroupId: groupId, Types: types}
			resp := &response.TotalResponse{}
			err := c.Call(context.Background(), "Total", req, resp)
			if err == nil {
				result <- resp.Total
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(result)
	for r := range result {
		total += r
	}
	return total
}


func ClientExit(uid int64) (result bool) {
	d := NewClient("Connection")
	services := d.GetServices()
	var wg sync.WaitGroup
	for _, service:= range services{
		ser := service
		go func() {
			wg.Add(1)
			d, _ := client.NewPeer2PeerDiscovery(ser.Key, "")
			c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.ExistRequest{Uid: uid}
			resp := &response.ExistResponse{}
			c.Call(context.Background(), "Exists", req, resp)
			if resp.Exists {
				result = true
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return result
}

