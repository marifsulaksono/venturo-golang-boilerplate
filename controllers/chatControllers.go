package controllers

import (
	"fmt"
	"net/http"
	"simple-crud-rnd/config"
	"simple-crud-rnd/helpers"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ChatController struct {
	db        *gorm.DB
	cfg       *config.Config
	wsManager *helpers.WebSocketManager
}

func NewChatController(db *gorm.DB, cfg *config.Config, wsManager *helpers.WebSocketManager) *ChatController {
	return &ChatController{db, cfg, wsManager}
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

	// Get user ID form query parameter for the client identifier
	userID := c.QueryParam("user_id")
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	client := &helpers.Client{
		ID:     userID,
		Socket: conn,
		Send:   make(chan helpers.Message),
	}

	// Register the client with the manager
	ch.wsManager.Register <- client

	// Handle incoming messages from the client
	go client.ListenForMessages(ch.wsManager)

	// Keep connection open
	select {}
}
