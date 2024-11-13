package models

import (
	"context"
	"fmt"
	"log"
	"simple-crud-rnd/helpers/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatModel struct {
	db *mongo.Database
}

func NewChatModel(db *mongo.Database) *ChatModel {
	return &ChatModel{
		db: db,
	}
}

const chatCollection = "chats"

func (cm *ChatModel) RetrieveChatHistory(userID, recipientID string) ([]utils.Message, error) {
	var history []utils.Message

	// Define a MongoDB filter to get messages where the user is either sender or receiver
	filter := bson.M{
		"$or": []bson.M{
			{"sender": userID, "receiver": recipientID},
			{"sender": recipientID, "receiver": userID},
		},
	}

	// Fetch messages from MongoDB
	cursor, err := cm.db.Collection(chatCollection).Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chat history: %w", err)
	}
	defer cursor.Close(context.Background())

	// Decode messages into the history slice
	if err := cursor.All(context.Background(), &history); err != nil {
		return nil, fmt.Errorf("failed to decode chat history: %w", err)
	}

	return history, nil
}

func (cm *ChatModel) InsertMessage(message utils.Message) error {
	return utils.InsertChatMessage(message, chatCollection)
}

func (cm *ChatModel) DeleteMessageByID(messageID string) error {
	objID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}

	_, err = cm.db.Collection(chatCollection).DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	log.Printf("Message %s deleted successfully", messageID)
	return nil
}
