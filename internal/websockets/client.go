package websockets

import "github.com/gorilla/websocket"

type Client struct {
	id      string
	conn    *websocket.Conn
	receive chan []byte
}

func NewClient(id string, conn *websocket.Conn) *Client {

	c := Client{
		id:      id,
		conn:    conn,
		receive: make(chan []byte),
	}

	return &c
}
