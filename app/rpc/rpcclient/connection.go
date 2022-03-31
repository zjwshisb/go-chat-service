package rpcclient

import (
	"context"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"sync"
	"ws/app/rpc/request"
	"ws/app/rpc/response"
)

func ConnectionTotal(groupId int64, types string) int64 {
	d := NewClient("Connection")
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
			req := &request.ConnectionGroupRequest{GroupId: groupId, Types: types}
			resp := &response.ConnectionTotalResponse{}
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

func ConnectionIds(groupId int64, types string) []int64 {
	d := NewClient("Connection")
	services := d.GetServices()
	ids := make([]int64, 0)
	result := make(chan []int64, len(services))
	var wg sync.WaitGroup
	fmt.Println("总服务")
	fmt.Println(len(services))
	for _, service := range services {
		ser := service
		fmt.Println(ser.Value)
		fmt.Printf("%v", ser)
		wg.Add(1)
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(ser.Key, "")
			c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.ConnectionGroupRequest{GroupId: groupId, Types: types}
			resp := &response.ConnectionIdsResponse{}
			err := c.Call(context.Background(), "Ids", req, resp)
			fmt.Println(resp.Ids)
			if err == nil {
				result <- resp.Ids
			} else {
				fmt.Println(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(result)
	for r := range result {
		fmt.Println("接受")
		fmt.Println(r)
		ids = append(ids, r...)
	}
	return ids
}

func ConnectionExist(uid int64) (result bool) {
	d := NewClient("Connection")
	services := d.GetServices()
	var wg sync.WaitGroup
	for _, service := range services {
		ser := service
		wg.Add(1)
		go func() {
			d, _ := client.NewPeer2PeerDiscovery(ser.Key, "")
			c := client.NewXClient("Connection", client.Failtry, client.RandomSelect, d, client.DefaultOption)
			defer c.Close()
			req := &request.ConnectionExistRequest{Uid: uid}
			resp := &response.ConnectionExistResponse{}
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
