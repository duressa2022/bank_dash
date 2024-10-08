package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"working.com/bank_dash/internal/domain"
	"working.com/bank_dash/package/mongo"
)

// type for working on the chat repository
type ChatRepository struct {
	database   mongo.Database
	collection string
}

const (
	Limit = 10
)

// method for creating new chat repository
func NewChatRepository(db mongo.Database, collection string) *ChatRepository {
	return &ChatRepository{
		database:   db,
		collection: collection,
	}
}

// method for storing chat message into the database
func (cr *ChatRepository) StoreMessage(c context.Context, message *domain.ChatMessage) (*domain.ChatResponse, error) {
	collection := cr.database.Collection(cr.collection)
	chatId, err := collection.InsertOne(c, message)
	if err != nil {
		return nil, err
	}

	var chatData domain.ChatResponse
	err = collection.FindOne(c, bson.D{{Key: "_id", Value: chatId}}).Decode(&chatData)
	if err != nil {
		return nil, err
	}

	return &chatData, nil
}

// method getting chat history from the database based id
func (cr *ChatRepository) GetMessage(c context.Context, id string) ([]*domain.ChatResponse, error) {
	cr.DeleteChatMessage(c, id, Limit)

	collection := cr.database.Collection(cr.collection)
	UserId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(c, bson.D{{Key: "_userId", Value: UserId}})
	if err != nil {
		return nil, err
	}

	var chatHistroy []*domain.ChatResponse
	for cursor.Next(c) {
		var histroy *domain.ChatResponse
		err := cursor.Decode(&histroy)
		if err != nil {
			return nil, err
		}
		chatHistroy = append(chatHistroy, histroy)
	}
	return chatHistroy, nil
}

// method for deleting the message if the limit is reached
func (cr *ChatRepository) DeleteChatMessage(c context.Context, id string, limit int64) error {
	collection := cr.database.Collection(cr.collection)
	userid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	chatCounter, err := collection.CountDocuments(c, bson.D{{Key: "_userId", Value: userid}})
	if err != nil {
		return err
	}
	if chatCounter < limit {
		return nil
	}

	opts := options.Find().SetSort(bson.D{{Key: "timeStamp", Value: 1}}).SetLimit(chatCounter - limit)
	cursor, err := collection.Find(c, bson.D{{Key: "_userId", Value: userid}}, opts)
	if err != nil {
		return err
	}

	var deletedChat []primitive.ObjectID
	for cursor.Next(c) {
		var chatMessage domain.ChatMessage
		err := cursor.Decode(&chatMessage)
		if err != nil {
			return err
		}
		deletedChat = append(deletedChat, chatMessage.ID)
	}

	for _, id := range deletedChat {
		_, err := collection.DeleteOne(c, bson.D{{Key: "_id", Value: id}})
		if err != nil {
			return err
		}
	}
	return nil

}
