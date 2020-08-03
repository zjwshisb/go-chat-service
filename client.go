package main

import (
	"errors"
	"github.com/gorilla/websocket"
)
var id = 1
type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Id int
}

func NewClient(conn *websocket.Conn) *Client {
	id ++
	return &Client{
		Conn: conn,
		Send: make(chan []byte, 64),
		Id: id,
	}
}
func (client *Client) readMessage() error  {
	return errors.New("")
}