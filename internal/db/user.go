package db

import (
	"context"
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"switchboard/internal/util"
	"time"

	"github.com/graph-gophers/dataloader"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUser(user *models.CreateUserRequest) (*models.User, *common.DetailedError) {
	hashedPassword, err := common.CreateHash(user.Password)
	if err != nil {
		return nil, GetDbError(err)
	}

	currentTime := time.Now()
	userId := util.UUIDv4()
	newUser := &models.User{
		ID:        userId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  string(hashedPassword),
		CreatedAt: &currentTime,
		UpdatedAt: &currentTime,
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
		if findError == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, GetDbError(fmt.Errorf("could not retrieve created document %s", userId))
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
		return nil, GetDbError(err)
	}
	return &user, nil
}

func BatchLoadUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	results := make([]*dataloader.Result, len(keys))
	usersCol := Database.Collection(USERS_COLLECTION)
	dbQuery := bson.D{
		{Key: "id", Value: bson.D{{
			Key: "$in", Value: keys,
		}}},
	}

	cursor, errFind := usersCol.Find(ctx, dbQuery)
	if errFind != nil {
		return []*dataloader.Result{{
			Data:  nil,
			Error: errFind,
		}}
	}
	users := make([]models.User, 0)
	err := cursor.All(ctx, &users)
	if err != nil {
		return []*dataloader.Result{{
			Data:  nil,
			Error: errFind,
		}}
	}

	for i := 0; i < len(keys); i++ {
		results[i] = &dataloader.Result{}
		for _, s := range users {
			if s.ID == keys[i].String() {
				results[i] = &dataloader.Result{
					Data:  &s,
					Error: nil,
				}
				break
			}
		}
	}

	return results
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
		return []models.User{}, GetDbError(errFind)
	}
	result := make([]models.User, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, GetDbError(err)
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
		return nil, GetDbError(err)
	}
	passwordVerified, err := common.VerifyHash(password, []byte(user.Password))
	if err != nil {
		return nil, GetDbError(err)
	}
	if passwordVerified {
		return &user, nil
	}
	return nil, nil
}
