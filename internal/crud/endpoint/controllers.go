package endpoint

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/common/db"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateEndpoint(endpoint Endpoint) (*Endpoint, *common.DetailedError) {
	eId, _ := uuid.NewRandom()
	endpointId := eId.String()
	currentTime := time.Now()
	newEndpoint := &Endpoint{
		ID:            endpointId,
		Path:          endpoint.Path,
		Method:        endpoint.Method,
		ResponseDelay: endpoint.ResponseDelay,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endpointsCollection := db.Database.Collection("endpoints")
	_, insertErr := endpointsCollection.InsertOne(ctx, newEndpoint)
	if insertErr != nil {
		return nil, db.WrapDBErrorIfNecessary(insertErr)
	}
	var createdEndpoint Endpoint
	findErr := endpointsCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: endpointId,
	}}).Decode(&createdEndpoint)
	if findErr != nil {
		return nil, common.WrapAsDetailedError(findErr)
	}
	return &createdEndpoint, nil
}
