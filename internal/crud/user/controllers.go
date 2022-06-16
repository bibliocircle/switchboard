package user

import (
	"context"
	"fmt"
	"switchboard/internal/common/db"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/common/security"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(user *CreateUserRequest) (*User, *err_utils.DetailedError) {
	userId, err := uuid.NewRandom()
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}

	hashedPassword, err := security.CreateHash(user.Password)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}

	currentTime := time.Now()
	newUser := &User{
		ID:        userId.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  string(hashedPassword),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userCollection := db.Database.Collection("users")
	_, insertError := userCollection.InsertOne(ctx, newUser)
	if insertError != nil {
		return nil, db.GetDbError(insertError)
	}

	var createdUser User
	findError := userCollection.FindOne(ctx, bson.D{{
		Key: "id", Value: userId.String(),
	}}).Decode(&createdUser)
	if findError != nil {
		return nil, err_utils.WrapAsDetailedError(fmt.Errorf("could not retrieve created document %s", userId))
	}
	return &createdUser, nil
}

func GetUserByID(userId string) (*User, *err_utils.DetailedError) {
	var user User
	userCollection := db.Database.Collection("users")
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
		return nil, err_utils.WrapAsDetailedError(err)
	}
	return &user, nil
}

func GetUserByEmailPassword(email string, password string) (*User, *err_utils.DetailedError) {
	var user User
	userCollection := db.Database.Collection("users")
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
		return nil, err_utils.WrapAsDetailedError(err)
	}
	passwordVerified, err := security.VerifyHash(password, []byte(user.Password))
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	if passwordVerified {
		return &user, nil
	}
	return nil, nil
}
