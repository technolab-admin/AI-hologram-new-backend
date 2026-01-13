package websockets


import (
	"log"
	"fmt"
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
			log.Printf("Client %v Error: %v\n", clientId, err)
			continue
		}

		if err = vallidateMsg(msg); err != nil {
			log.Printf("Client %v Error: Validating message: %v\n", clientId, err)
			continue
		}

		targetId := msg["target"]
		targetClient, targetIsActive := s.clients[targetId]
		if !targetIsActive {
			log.Printf("Client %v Error sending message to %v, target not connected\n", clientId, targetId)
			continue
		}

		err = targetClient.conn.WriteMessage(websocket.TextMessage, rawJSON)

		if err != nil {
			log.Printf("Client %v Error sending message to %v, closing connection: %v\n", clientId, targetId, err)

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

	if !isKnownClient(clientId) {
		log.Printf("Error: Client %v failed to connect: Client is unknown\n", clientId)
		conn.Close()
		return
	}

	client := NewClient(clientId, conn)
	log.Printf("Client %v connected\n", client.id)

	s.mutex.Lock()
	s.clients[client.id] = client
	s.mutex.Unlock()

	go s.listener(client.id)
	go s.messageHandler(client.id)
}

func (s *Server) Start() error {

	log.Printf("WebSocket server started on port %v\n", s.port)

	err := http.ListenAndServe(s.port, nil)
	if err != nil {
		return fmt.Errorf("Error starting server: %v\n", err)
	}
	return nil
}
