package db

import (
	"context"
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"switchboard/internal/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUser(user *models.CreateUserRequest) (*models.User, *common.DetailedError) {
	hashedPassword, err := common.CreateHash(user.Password)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}

	currentTime := time.Now()
	userId := util.UUIDv4()
	newUser := &models.User{
		ID:        userId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  string(hashedPassword),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := Database.Collection(USERS_COLLECTION)
	_, insertError := userCollection.InsertOne(ctx, newUser)
	if insertError != nil {
		return nil, GetDbError(insertError)
	}

	var createdUser models.User
	findError := userCollection.FindOne(ctx, bson.D{{
		Key: "id", Value: userId,
	}}).Decode(&createdUser)
	if findError != nil {
		return nil, common.WrapAsDetailedError(fmt.Errorf("could not retrieve created document %s", userId))
	}
	return &createdUser, nil
}

func GetUserByID(userId string) (*models.User, *common.DetailedError) {
	var user models.User
	userCollection := Database.Collection(USERS_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := userCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: userId,
	}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, common.WrapAsDetailedError(err)
	}
	return &user, nil
}

func GetUsers() ([]models.User, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := Database.Collection(USERS_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}
	cursor, errFind := userCollection.Find(ctx, bson.D{}, findOpts)
	if errFind != nil {
		return []models.User{}, common.WrapAsDetailedError(errFind)
	}
	result := make([]models.User, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return result, nil
}

func GetUserByEmailPassword(email string, password string) (*models.User, *common.DetailedError) {
	var user models.User
	userCollection := Database.Collection(USERS_COLLECTION)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := userCollection.FindOne(ctx, bson.D{{
		Key:   "email",
		Value: email,
	}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, common.WrapAsDetailedError(err)
	}
	passwordVerified, err := common.VerifyHash(password, []byte(user.Password))
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	if passwordVerified {
		return &user, nil
	}
	return nil, nil
}
