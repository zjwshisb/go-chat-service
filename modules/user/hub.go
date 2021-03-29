package user

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"ws/db"
)

type hub struct {
	Clients map[int64]*UClient
	Login chan *UClient
	Logout chan *UClient
	Lock sync.Mutex
}
func (hub *hub) run()  {
	for {
		select {
		case client := <- hub.Login:
			hub.Lock.Lock()
			hub.Clients[client.UserId] = client
			ctx := context.Background()
			_, err := db.Redis.Get(ctx, fmt.Sprintf("user:%d:server", client.UserId)).Result()
			if err == nil {
				db.Redis.SAdd(ctx, "user:waiting:list", time.Now().Unix(), client.UserId)
			}
			hub.Lock.Unlock()
		case client := <- hub.Logout:
			hub.Lock.Lock()
			delete(hub.Clients, client.UserId)
			log.Print("logout")
			hub.Lock.Unlock()
		}
	}
}

