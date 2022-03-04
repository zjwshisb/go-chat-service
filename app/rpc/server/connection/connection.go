package connection

import (
	"context"
	"github.com/golang-module/carbon"
	"ws/app/websocket"
)

type Connection struct {

}

type Reply struct {
	Data []map[string]string
}

type Args struct {
	types string
}

func (conn *Connection) Count(ctx context.Context, args *Args, reply *Reply)  {

}

func (conn *Connection) Show(ctx context.Context, args *Args, reply *Reply) error {
	conns := websocket.AdminManager.GetTotalConn()
	reply.Data = make([]map[string]string, 0)
	for _, conn := range conns {
		reply.Data = append(reply.Data, map[string]string{
			"uid" : conn.GetUid(),
			"created_at": carbon.CreateFromTimestamp(conn.GetCreateTime()).ToDateTimeString(),
		})
	}
	return nil
}
