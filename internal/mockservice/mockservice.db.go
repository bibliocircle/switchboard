package mockservice

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateMockService(userId string, ms *CreateMockServiceRequestBody) (*MockService, *common.DetailedError) {
	currentTime := time.Now()
	newMockService := &MockService{
		ID:   ms.ID,
		Name: ms.Name,
		Type: ms.Type,
		Config: GlobalMockServiceConfig{
			InjectHeaders:       ms.Config.InjectHeaders,
			GlobalResponseDelay: ms.Config.GlobalResponseDelay,
		},
		CreatedBy: userId,
		CreatedAt: &currentTime,
		UpdatedAt: &currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCollection := db.Database.Collection(db.MOCKSERVICES_COLLECTION)
	_, insertErr := mockServicesCollection.InsertOne(ctx, newMockService)
	if insertErr != nil {
		return nil, db.GetDbError(insertErr)
	}
	var createdMockService MockService
	findErr := mockServicesCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: ms.ID,
	}}).Decode(&createdMockService)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, db.GetDbError(findErr)
	}
	return &createdMockService, nil
}

func GetMockServices() (*[]*MockService, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCol := db.Database.Collection(db.MOCKSERVICES_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}

	cursor, errFind := mockServicesCol.Find(ctx, bson.D{}, findOpts)
	if errFind != nil {
		return &[]*MockService{}, db.GetDbError(errFind)
	}
	result := make([]*MockService, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return &[]*MockService{}, db.GetDbError(err)
	}
	return &result, nil
}

func GetMockServicesByIds(ids []string) (*[]*MockService, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCol := db.Database.Collection(db.MOCKSERVICES_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}

	cursor, errFind := mockServicesCol.Find(ctx, bson.D{
		{Key: "id", Value: bson.D{{Key: "$in", Value: ids}}},
	}, findOpts)
	if errFind != nil {
		return &[]*MockService{}, db.GetDbError(errFind)
	}
	result := make([]*MockService, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, db.GetDbError(err)
	}
	return &result, nil
}

func GetMockServiceByID(ID string) (*MockService, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var ms MockService
	mockServicesCol := db.Database.Collection(db.MOCKSERVICES_COLLECTION)
	findErr := mockServicesCol.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: ID,
	}}).Decode(&ms)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, db.GetDbError(findErr)
	}
	return &ms, nil
}

func DeleteMockService(userID, mockServiceID string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCol := db.Database.Collection(db.MOCKSERVICES_COLLECTION)
	result, errDel := mockServicesCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: mockServiceID},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, db.GetDbError(errDel)
	}
	return result.DeletedCount > 0, nil
}
