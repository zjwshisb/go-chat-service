package uclient

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type Client struct {
	Conn       *websocket.Conn
	Id         int64
	intMessage chan []byte
	outMessage chan []byte
	isClose    bool
	flag       chan bool
	lock       sync.Mutex
}

var (
	WaitAccepts = make(map[*Client]bool)
)

func NewClient(conn *websocket.Conn, id int64) *Client {
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
	logrus.Info("client:" + strconv.FormatInt(client.Id, 10) + ":connect")
	return client
}
func (client *Client) isCLose() bool {
	defer client.lock.Unlock()
	client.lock.Lock()
	return client.isClose
}
func (client *Client) safeClose() {
	defer client.lock.Unlock()
	client.lock.Lock()
	if !client.isClose {
		client.isClose = true
		err := client.Conn.Close()
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("client:" + strconv.FormatInt(client.Id, 10) + ":close")
		close(client.flag)
	}
}

func (client *Client) Push(msg []byte) {
	client.outMessage <- msg
}

func (client *Client) safeSend() {
	defer client.safeClose()
	LOOP: for  {
		select {
			case <-client.flag:
				break
			default:
				msg := <-client.outMessage
				err := client.Conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					logrus.Error("client:" + strconv.FormatInt(client.Id, 10) + ":break send loop:" + err.Error())
					break LOOP
				}
				logrus.Info("client:" + strconv.FormatInt(client.Id, 10) + ":send:" + string(msg))

		}
	}
}
func (client *Client) SafeRead() {
	defer client.safeClose()
	LOOP: for {
		select {
			case <-client.flag:
				break
			default:
				_, msg, err := client.Conn.ReadMessage()
				if err != nil {
					logrus.Error("client:" + strconv.FormatInt(client.Id, 10) + ":break read loop:" + err.Error())
					break LOOP
				}
				client.intMessage <- msg
				logrus.Error("client:" + strconv.FormatInt(client.Id, 10) + ":receive:" + string(msg))
				client.Push(msg)
		}
	}
}
