package websockets

type WSConn interface { // This replaces *websocket.Conn to make testing more scalable
    ReadMessage() (int, []byte, error)
    WriteMessage(int, []byte) error
    Close() error
}

type Client struct {
	id      string
	conn    WSConn
	receive chan []byte
}

func NewClient(id string, conn WSConn) *Client {

	c := Client{
		id:      id,
		conn:    conn,
		receive: make(chan []byte),
	}

	return &c
}
