package uclient

import (
	"github.com/gorilla/websocket"
	"sync"
	"fmt"
)

type Client struct {
	Conn       *websocket.Conn
	Id         int32
	intMessage chan []byte
	outMessage chan []byte
	isClose    bool
	flag       chan bool
	lock       sync.Mutex
}

var (
	WaitAccepts = make(map[*Client]bool)
)

func NewClient(conn *websocket.Conn, id int32) *Client {
	client := &Client{
		Conn:       conn,
		Id:         id,
		intMessage: make(chan []byte, 1000),
		outMessage: make(chan []byte, 1000),
		isClose:    false,
		flag: make(chan bool),
	}
	go client.safeSend()
	go client.SafeRead()
	return client
}
func (client *Client) isCLose() bool {
	defer client.lock.Unlock()
	client.lock.Lock()
	return client.isClose
}
func (client *Client) safeClose() {
	client.lock.Lock()
	if !client.isClose {
		client.isClose = true
		err := client.Conn.Close()
		if err != nil {
			fmt.Println(err)
		}
		close(client.flag)
		fmt.Println("client:" + string(client.Id) + " cloase")
	}
	client.lock.Unlock()
}

func (client *Client) Push(msg []byte) {
	client.outMessage <- msg
}

func (client *Client) safeSend() {
	defer client.safeClose()
	for  {
		select {
			case <-client.flag:
				break
			default:
				msg := <-client.outMessage
				err := client.Conn.WriteMessage(websocket.TextMessage, msg)
				if err == nil {
					break
				}
		}
	}
}
func (client *Client) SafeRead() {
	defer client.safeClose()
	for {
		select {
			case <-client.flag:
				break
			default:
				_, msg, err := client.Conn.ReadMessage()
				if err == nil {
					return
				}
				fmt.Println(msg)
				client.intMessage <- msg
				client.Push(msg)
		}
	}
}
