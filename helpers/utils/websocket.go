package utils

import (
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

/*
WebSocketManager handles WebSocket connections and message broadcasting

Clients: a map of user ID to WebSocket connections
Broadcast: a channel to broadcast messages
Register: a channel to register new clients
Unregister: a channel to unregister clients

@Authored by Muhammad Arif Sulaksono
*/
type WebSocketManager struct {
	Clients    map[string]*websocket.Conn
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	Sender   string    `json:"sender"`
	Receiver string    `json:"receiver"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[string]*websocket.Conn),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (manager *WebSocketManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client.ID] = client.Socket
			log.Printf("Client %s connected", client.ID)

		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client.ID]; ok {
				close(client.Send)
				delete(manager.Clients, client.ID)
				log.Printf("Client %s successfully disconnected", client.ID)
			}

		case message := <-manager.Broadcast:
			chatCollection := "chats"
			InsertChatMessage(message, chatCollection)

			// send messages to the recipient (personal chat for now)
			for clientID, conn := range manager.Clients {
				if strings.EqualFold(clientID, message.Receiver) {
					err := conn.WriteJSON(message)
					if err != nil {
						log.Printf("Error sending message to client %s: %v", clientID, err)
						conn.Close()
						delete(manager.Clients, clientID)
					} else {
						log.Printf("Message sent to %s", clientID)
					}
				}
			}
		}
	}
}

/*
Client represents a single WebSocket connection
*/
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan Message
}

func (client *Client) ListenForMessages(manager *WebSocketManager) {
	defer func() {
		manager.Unregister <- client
		client.Socket.Close()
	}()

	// Keep receiving messages from the WebSocket
	for {
		var msg Message
		err := client.Socket.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Client %s disconnected gracefully", client.ID)
			} else {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		msg.Time = time.Now().Truncate(time.Second) // Truncate(time.Second) removes any sub-second precision

		log.Printf("[%s] Received message from %s to %s: %s", msg.Time, msg.Sender, msg.Receiver, msg.Content)
		// Send message to broadcast channel
		manager.Broadcast <- msg
	}
}
