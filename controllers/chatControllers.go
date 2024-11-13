package controllers

import (
	"fmt"
	"log"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"
	"simple-crud-rnd/helpers/utils"
	"simple-crud-rnd/models"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatController struct {
	mongo     *mongo.Database
	chatModel *models.ChatModel
	cfg       *config.Config
	wsManager *utils.WebSocketManager
}

func NewChatController(mongo *mongo.Database, chatModel *models.ChatModel, cfg *config.Config, wsManager *utils.WebSocketManager) *ChatController {
	return &ChatController{mongo, chatModel, cfg, wsManager}
}

/*
*
* WebSocketHandler handles WebSocket connections and message broadcasts
* This function sets up a WebSocket connection, registers the client with the manager, and starts a goroutine to handle incoming messages from the client.
*
* @Authored by Muhammad Arif Sulaksono
*
 */
func (ch *ChatController) WebSocketHandler(c echo.Context) error {
	// Set up WebSocket connection
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // allow all origins for simplicity
	}).Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("could not upgrade to WebSocket: %v", err)
	}

	// Get user ID and recipient ID from query parameters
	// userID := c.QueryParam("user_id")
	recipientID := c.QueryParam("recipient_id")
	accessToken := c.QueryParam("access_token")

	user, err := helpers.VerifyTokenJWT(accessToken, false)
	if err != nil {
		log.Printf("Missing sender_id or recipient_id")
		conn.WriteJSON(map[string]string{"error": "Unauthorized, please login first"})
		conn.Close()
		return nil
	}

	// Validate user ID and recipient ID
	if recipientID == "" {
		log.Printf("Missing sender_id or recipient_id")
		conn.WriteJSON(map[string]string{"error": "sender_id and recipient_id are required"})
		conn.Close()
		return nil
	}

	// Fetch and send chat history to the user
	history, err := ch.chatModel.RetrieveChatHistory(user.Email, recipientID)
	if err != nil {
		log.Printf("Failed to retrieve chat history for user %s: %v", user.Email, err)
		conn.WriteJSON(map[string]string{"error": "failed to retrieve chat history"})
		conn.Close()
		return nil
	}

	for _, msg := range history {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending history message to client %s: %v", user.Email, err)
			conn.Close()
			return nil
		}
	}

	// Create a new client and register it
	client := &utils.Client{
		ID:     user.Email,
		Socket: conn,
		Send:   make(chan utils.Message),
	}

	ch.wsManager.Register <- client

	// Start a goroutine to listen for messages from the client
	go client.ListenForMessages(ch.wsManager)

	// Keep WebSocket connection open without returning an HTTP response
	return nil
}
