package client

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"sync"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

func NoticeRepeatConnect(id int64, types string, newUuid string, server string) {
	d, _ := client.NewPeer2PeerDiscovery(server, "")
	c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()
	req := &request.RepeatConnectRequest{Types: types, Id: id, NewUuid: newUuid}
	resp := &response.NilResponse{}
	_ = c.Call(context.Background(), "RepeatConnect", req, resp)
}

func ConnectionOnline(id int64, types string, server string) bool {
	d, _ := client.NewPeer2PeerDiscovery(server, "")
	c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()
	req := &request.OnlineRequest{Types: types, Id: id}
	resp := &response.OnlineResponse{}
	err := c.Call(context.Background(), "Online", req, resp)
	if err != nil {

	}
	return resp.Data
}

func ConnectionTotal(groupId int64, types string) int64 {
	d := NewDiscovery("Connection")
	services := d.GetServices()
	var total int64
	result := make(chan int64, len(services))
	var wg sync.WaitGroup
	for _, service := range services {
		wg.Add(1)
		ser := service
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(ser.Key, "")
			c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.NormalRequest{GroupId: groupId, Types: types}
			resp := &response.CountResponse{}
			err := c.Call(context.Background(), "Count", req, resp)
			if err == nil {
				result <- resp.Data
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

func ConnectionIds(groupId int64, types string) []int64 {
	d := NewDiscovery("Connection")
	services := d.GetServices()
	ids := make([]int64, 0)
	result := make(chan []int64, len(services))
	var wg sync.WaitGroup
	for _, service := range services {
		ser := service
		wg.Add(1)
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(ser.Key, "")
			c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.IdsRequest{GroupId: groupId, Types: types}
			resp := &response.IdsResponse{}
			err := c.Call(context.Background(), "Ids", req, resp)
			if err == nil {
				result <- resp.Data
			} else {
				fmt.Println(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(result)
	for r := range result {
		ids = append(ids, r...)
	}
	return ids
}

func ConnectionAllCount(types string) int64 {
	d := NewDiscovery("Connection")
	services := d.GetServices()
	var total int64
	result := make(chan int64, len(services))
	var wg sync.WaitGroup
	for _, service := range services {
		wg.Add(1)
		ser := service
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(ser.Key, "")
			c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.TypeRequest{Types: types}
			resp := &response.CountResponse{}
			err := c.Call(context.Background(), "AllCount", req, resp)
			if err == nil {
				result <- resp.Data
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
