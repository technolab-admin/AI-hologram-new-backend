package meshy

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"AI-HOLOGRAM-NEW-BACKEND/internal/config"
)

type WSClient struct {
	id   string
	send chan []byte
	conn *websocket.Conn
}

func NewWSClient(id string) *WSClient {
	return &WSClient{id: id, send: make(chan []byte, 16)}
}

func (wsc *WSClient) StartWebsocketClient() {

	for {
		cfg, err := config.Load()
		if err != nil {
			continue
		}

		websocketUrl := cfg.WebsocketURL
		url := fmt.Sprintf("ws://%v/ws?id=%v", websocketUrl, wsc.id)

		wsc.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			continue
		}

		wsc.frontendNotifier()
	}
}

func (wsc *WSClient) frontendNotifier() {

	for {
		msg := <-wsc.send

		err := wsc.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			wsc.conn.Close()
			log.Printf("Error sending message: %v", err) // Change to logger function
			return
		}
	}
}

func (wsc *WSClient) notifyFrontend(msg map[string]string) error {

	rawJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Error converting message to bytes: %v", err)
	}
	wsc.send <- rawJSON
	return nil
}
