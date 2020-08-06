package uclient

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
	Id int32
}
var (
	WaitAccepts map[*Client]bool
)
func NewClient(conn *websocket.Conn, id int32) *Client {
	return &Client{
		Conn: conn,
		Id: id,
	}
}
