package mockservice

import (
	"context"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GlobalMockServiceConfig struct {
	InjectHeaders map[string]string `json:"injectHeaders" bson:"injectHeaders,omitempty"`
}

type MockService struct {
	ID        string                  `json:"id" bson:"id,omitempty"`
	Name      string                  `json:"name" bson:"name,omitempty"`
	Type      string                  `json:"type" bson:"type,omitempty"`
	Config    GlobalMockServiceConfig `json:"config" bson:"config,omitempty"`
	CreatedBy string                  `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt time.Time               `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time               `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func CreateMockService(userId string, ms *CreateMockServiceRequestBody) (*MockService, *err_utils.DetailedError) {
	currentTime := time.Now()
	newMockService := &MockService{
		ID:   ms.ID,
		Name: ms.Name,
		Type: ms.Type,
		Config: GlobalMockServiceConfig{
			InjectHeaders: ms.Config.InjectHeaders,
		},
		CreatedBy: userId,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
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
		return nil, err_utils.WrapAsDetailedError(findErr)
	}
	return &createdMockService, nil
}

func GetMockServices() ([]MockService, *err_utils.DetailedError) {
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
		return []MockService{}, err_utils.WrapAsDetailedError(errFind)
	}
	result := make([]MockService, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	return result, nil
}

func DeleteMockService(userID, mockServiceID string) (bool, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mockServicesCol := db.Database.Collection(db.MOCKSERVICES_COLLECTION)
	result, errDel := mockServicesCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: mockServiceID},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, err_utils.WrapAsDetailedError(errDel)
	}
	return result.DeletedCount > 0, nil
}
