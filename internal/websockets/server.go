package websockets

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	port     string
	mutex    *sync.Mutex
	clients  map[string]*Client
	upgrader websocket.Upgrader
}


func (s *Server) Start() error {

	log.Printf("WebSocket server started on port %v\n", s.port)

	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		return fmt.Errorf("Error starting server: %v\n", err)
	}
	return nil
}


func (s *Server) MessageHandler(clientId string) {
    s.messageHandler(clientId)
}

func (s *Server) Listener(clientId string) {
    s.listener(clientId)
}

func (s *Server) AddClient(client *Client) {
	s.mutex.Lock()
	s.clients[client.id] = client
	s.mutex.Unlock()
}


func (s *Server) GetNClients() int {
	s.mutex.Lock()
	n := len(s.clients)
	s.mutex.Unlock()
	return n
}



func NewServer(port string) *Server {

	s := Server{
		port:    port,
		mutex:   &sync.Mutex{},
		clients: make(map[string]*Client),

		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	http.HandleFunc("/ws", s.wsHandler)

	return &s
}

func (s *Server) messageHandler(clientId string) {

	s.mutex.Lock()
	client := s.clients[clientId]
	s.mutex.Unlock()

	for {

		rawJSON := <-client.receive

		msg, err := UnmarshalMessage(rawJSON)
		if err != nil {
			log.Printf("Client %v Error: %v\n", clientId, err) // Change to logger function
			continue
		}

		if err = VallidateMsg(msg); err != nil {
			log.Printf("Client %v Error: Validating message: %v\n", clientId, err) // Change to logger function
			continue
		}

		targetId := msg["target"]
		targetClient, targetIsActive := s.clients[targetId]
		if !targetIsActive {
			log.Printf("Client %v Error sending message to %v, target not connected\n", clientId, targetId) // Change to logger function
			continue
		}

		err = targetClient.conn.WriteMessage(websocket.TextMessage, rawJSON)

		if err != nil {
			log.Printf("Client %v Error sending message to %v, closing connection: %v\n", clientId, targetId, err) // Change to logger function

			s.mutex.Lock()
			targetClient.conn.Close()
			delete(s.clients, targetId)
			s.mutex.Unlock()
			break
		}
	}
}

func (s *Server) listener(clientId string) {

	s.mutex.Lock()
	client := s.clients[clientId]
	s.mutex.Unlock()

	for {
		_, rawJSON, err := client.conn.ReadMessage()

		if err != nil {
			log.Printf("Client %v disconected: %v\n", clientId, err)

			s.mutex.Lock()
			client.conn.Close()
			delete(s.clients, clientId)
			s.mutex.Unlock()
			break
		}

		client.receive <- rawJSON
	}
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading: %v\n", err)
		return
	}

	clientId := r.URL.Query().Get("id")

	if !IsKnownClient(clientId) {
		log.Printf("Error: Client %v failed to connect: Client is unknown\n", clientId)
		conn.Close()
		return
	}

	client := NewClient(clientId, conn)
	log.Printf("Client %v connected\n", client.id)

	s.AddClient(client)

	go s.listener(client.id)
	go s.messageHandler(client.id)
}


