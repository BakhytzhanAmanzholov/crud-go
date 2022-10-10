package repositories

import (
	"context"
	"crud-golang/internal/models"
	"crud-golang/pkg/client/mongodb"
	"crud-golang/pkg/logging"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var accountsCollections *mongo.Collection = mongodb.GetCollection(mongodb.DB, "accounts")

type repository struct {
	accountsCollections *mongo.Collection
}

func (r repository) Create(user *models.Account) (models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := accountsCollections.InsertOne(ctx, user)
	if err != nil {
		logging.Error(err)
		return models.Account{}, errors.New("database error")
	}

	logging.Infof("Create account by email", user.Email)
	return *user, err
}

func (r repository) FindAll() (u []models.Account, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var accounts []models.Account
	defer cancel()

	results, err := accountsCollections.Find(ctx, bson.M{})

	if err != nil {
		logging.Fatal(err)
		return nil, errors.New("database error")
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var account models.Account
		if err = results.Decode(&account); err != nil {
			logging.Fatal(err)
			return nil, errors.New("decode error")
		}

		accounts = append(accounts, account)
	}
	logging.Infof("find all accounts")
	return accounts, err
}

func (r repository) FindOne(id string) (models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.Account
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	err := accountsCollections.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

	if err != nil {
		logging.Error(err)
		return models.Account{}, errors.New("database error")
	}

	logging.Infof("Find account by id, email: ", user.Email)
	return user, err
}

func (r repository) Update(id string, user models.Account) (models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	update := bson.M{"email": user.Email, "password": user.Password, "username": user.Username}

	result, err := accountsCollections.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

	if err != nil {
		logging.Error(err)
		return models.Account{}, errors.New("database error")
	}

	var updatedUser models.Account
	if result.MatchedCount == 1 {
		err := accountsCollections.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			logging.Error(err)
			return models.Account{}, errors.New("database error")
		}
	}

	logging.Infof("Update account", user.Email)
	return updatedUser, err
}

func (r repository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	result, err := accountsCollections.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		logging.Error(err)
		return errors.New("database error")
	}

	if result.DeletedCount < 1 {

		return errors.New("incorrect id")
	}

	logging.Infof("Delete account by id", id)
	return nil
}

func (r repository) FindByEmail(email string) (models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"email": email}
	result := models.Account{}
	err := accountsCollections.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Account{}, nil
		}

		return models.Account{}, errors.New("incorrect email")
	}

	logging.Infof("Find account by email", email)
	return result, nil
}

func NewRepository() Repository {
	return &repository{}
}
