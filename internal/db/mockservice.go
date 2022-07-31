package db

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMockService(userId string, ms *models.CreateMockServiceRequestBody) (*models.MockService, *common.DetailedError) {
	currentTime := time.Now()
	newMockService := &models.MockService{
		ID:   ms.ID,
		Name: ms.Name,
		Type: ms.Type,
		Config: models.GlobalMockServiceConfig{
			InjectHeaders:       ms.Config.InjectHeaders,
			GlobalResponseDelay: ms.Config.GlobalResponseDelay,
		},
		CreatedBy: userId,
		CreatedAt: &currentTime,
		UpdatedAt: &currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCollection := Database.Collection(MOCKSERVICES_COLLECTION)
	_, insertErr := mockServicesCollection.InsertOne(ctx, newMockService)
	if insertErr != nil {
		return nil, GetDbError(insertErr)
	}
	var createdMockService models.MockService
	findErr := mockServicesCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: ms.ID,
	}}).Decode(&createdMockService)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, GetDbError(findErr)
	}
	return &createdMockService, nil
}

func GetMockServices() (*[]*models.MockService, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCol := Database.Collection(MOCKSERVICES_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}

	cursor, errFind := mockServicesCol.Find(ctx, bson.D{}, findOpts)
	if errFind != nil {
		return &[]*models.MockService{}, GetDbError(errFind)
	}
	result := make([]*models.MockService, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, GetDbError(err)
	}
	return &result, nil
}

func GetMockServicesByIds(ids []string) (*[]*models.MockService, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCol := Database.Collection(MOCKSERVICES_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}

	cursor, errFind := mockServicesCol.Find(ctx, bson.D{
		{Key: "id", Value: bson.D{{Key: "$in", Value: ids}}},
	}, findOpts)
	if errFind != nil {
		return &[]*models.MockService{}, GetDbError(errFind)
	}
	result := make([]*models.MockService, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, GetDbError(err)
	}
	return &result, nil
}

func GetMockServiceByID(ID string) (*models.MockService, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var ms models.MockService
	mockServicesCol := Database.Collection(MOCKSERVICES_COLLECTION)
	findErr := mockServicesCol.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: ID,
	}}).Decode(&ms)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, GetDbError(findErr)
	}
	return &ms, nil
}

func DeleteMockService(userID, mockServiceID string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCol := Database.Collection(MOCKSERVICES_COLLECTION)
	result, errDel := mockServicesCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: mockServiceID},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, GetDbError(errDel)
	}
	return result.DeletedCount > 0, nil
}
