package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ChatCollection = "chats"
)

// type for working with chatting domain
type ChatMessage struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"userid" bson:"_userId"`
	Message   string             `json:"message" bson:"message"`
	Response  string             `json:"response" bson:"response"`
	TimeStamp time.Time          `json:"timeStamp" bson:"timeStamp"`
}

// type for working with chat response
type ChatResponse struct {
	Message  string `json:"message" bson:"message"`
	Response string `json:"response" bson:"response"`
}

// type for working with chat request
type ChatRequest struct {
	Message string `json:"message" bson:"message"`
}

// interface for working with chat message repository
type ChatRepository interface {
	StoreMessage(c context.Context, userID string, message ChatMessage) (*ChatMessage, error)
	GetMessage(c context.Context, userId string) ([]*ChatMessage, error)
	DeleteChatMessage(c context.Context, userID string, messageLimit int64) error
}
