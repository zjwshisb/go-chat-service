package uclient

import "github.com/gorilla/websocket"

type Client struct {
	Conn *websocket.Conn
}
var (
	WaitAccepts map[*Client]bool
)
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}
