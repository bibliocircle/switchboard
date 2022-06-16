package mockservice

import (
	"context"
	"switchboard/internal/common/db"
	"switchboard/internal/common/err_utils"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type GlobalMockServiceConfig struct {
	InjectHeaders map[string]string `json:"injectHeaders" bson:"injectHeaders,omitempty"`
}

type MockService struct {
	ID        string                  `json:"id" bson:"id,omitempty"`
	Name      string                  `json:"name" bson:"name,omitempty"`
	Key       string                  `json:"key" bson:"key,omitempty"`
	Type      string                  `json:"type" bson:"type,omitempty"`
	Config    GlobalMockServiceConfig `json:"config" bson:"config,omitempty"`
	CreatedBy string                  `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt time.Time               `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time               `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func CreateMockService(userId string, endpoint *MockService) (*MockService, *err_utils.DetailedError) {
	eId, _ := uuid.NewRandom()
	mockSvcId := eId.String()
	currentTime := time.Now()
	newMockService := &MockService{
		ID:   mockSvcId,
		Name: endpoint.Name,
		Type: endpoint.Type,
		Config: GlobalMockServiceConfig{
			InjectHeaders: endpoint.Config.InjectHeaders,
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
		Value: mockSvcId,
	}}).Decode(&createdMockService)
	if findErr != nil {
		return nil, err_utils.WrapAsDetailedError(findErr)
	}
	return &createdMockService, nil
}
