package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"working.com/bank_dash/internal/domain"
	"working.com/bank_dash/package/mongo"
)

// type for working with bank service repos
type BankRepository struct {
	database   mongo.Database
	collection string
}

// method for working bank service repos
func NewBankRepository(db mongo.Database, collection string) *BankRepository {
	return &BankRepository{
		database:   db,
		collection: collection,
	}
}

// method for getting banks by using id
func (br *BankRepository) GetBankById(c context.Context, id string) (*domain.BankService, error) {
	collection := br.database.Collection(br.collection)
	company_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var Bank *domain.BankService
	err = collection.FindOne(c, bson.D{{Key: "_id", Value: company_id}}).Decode(&Bank)
	if err != nil {
		return nil, err
	}
	return Bank, err
}

// method for updating the bank information based on id
func (br *BankRepository) UpdateBank(c context.Context, id string, bankRequest *domain.BankRequest) (*domain.BankService, error) {
	collection := br.database.Collection(br.collection)
	UpdatingBank := bson.M{
		"name":          bankRequest.Name,
		"details":       bankRequest.Details,
		"numberOfUsers": bankRequest.NumberOfUsers,
		"status":        bankRequest.Status,
		"type":          bankRequest.Type,
		"icon":          bankRequest.Icon,
	}

	BankId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	updateResult, err := collection.UpdateOne(c, bson.D{{Key: "_id", Value: BankId}}, bson.M{"$set": UpdatingBank})
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("no matched docs are found")
	}
	if updateResult.ModifiedCount == 0 {
		return nil, errors.New("no modified docs are found")
	}

	var Bank *domain.BankService
	err = collection.FindOne(c, bson.D{{Key: "_id", Value: BankId}}).Decode(&Bank)
	if err != nil {
		return &domain.BankService{}, err
	}
	return Bank, nil
}

// method for deleting bank service from the database
func (br *BankRepository) DeleteBank(c context.Context, id string) error {
	collection := br.database.Collection(br.collection)
	BankId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(c, bson.D{{Key: "_id", Value: BankId}})
	return err
}

// method for posting the bankservice into the database
func (br *BankRepository) PostBank(c context.Context, bank *domain.BankService) (*domain.BankService, error) {
	collection := br.database.Collection(br.collection)

	var Existing *domain.BankService
	err := collection.FindOne(c, bson.D{{Key: "name", Value: bank.Name}}).Decode(&Existing)
	if err == nil {
		return nil, errors.New("already existing bank")
	}
	bankId, err := collection.InsertOne(c, bank)
	if err != nil {
		return nil, err
	}

	var BankService domain.BankService
	err = collection.FindOne(c, bson.D{{Key: "_id", Value: bankId}}).Decode(&BankService)
	if err != nil {
		return nil, err
	}
	return &BankService, nil
}

// method for searching by using name of the bank service
func (br *BankRepository) SearchByName(c context.Context, term string) (*domain.BankService, error) {
	collection := br.database.Collection(br.collection)
	var Bank *domain.BankService
	err := collection.FindOne(c, bson.D{{Key: "name", Value: term}}).Decode(&Bank)
	if err != nil {
		return nil, err
	}
	return Bank, nil
}

// method for getting companies based on the page and size
func (br *BankRepository) GetBanks(c context.Context, page int, size int) ([]*domain.BankService, int, error) {
	var Banks []*domain.BankService
	collection := br.database.Collection(br.collection)

	skip := (page - 1) * size
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(size))
	cursor, err := collection.Find(c, bson.D{{}}, opts)
	if err != nil {
		return nil, 0, err
	}
	for cursor.Next(c) {
		var bank *domain.BankService
		err := cursor.Decode(&bank)
		if err != nil {
			return nil, 0, err
		}
		Banks = append(Banks, bank)
	}

	totalnumber, err := collection.CountDocuments(c, bson.D{{}})
	if err != nil {
		return nil, 0, err
	}
	return Banks, int(totalnumber), nil

}
