package db

import (
	"context"
	"strings"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"switchboard/internal/util"
	"time"

	"github.com/graph-gophers/dataloader"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateEndpoint(userId string, ep *models.CreateEndpointRequestBody) (*models.Endpoint, *common.DetailedError) {
	endpointId := util.UUIDv4()
	currentTime := time.Now()
	newEndpoint := &models.Endpoint{
		ID:            endpointId,
		MockServiceId: ep.MockServiceId,
		Path:          ep.Path,
		Method:        strings.ToUpper(ep.Method),
		Description:   ep.Description,
		ResponseDelay: ep.ResponseDelay,
		CreatedBy:     userId,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endpointsCollection := Database.Collection(ENDPOINT_COLLECTION)
	_, insertErr := endpointsCollection.InsertOne(ctx, newEndpoint)
	if insertErr != nil {
		return nil, GetDbError(insertErr)
	}
	var createdEndpoint models.Endpoint
	findErr := endpointsCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: endpointId,
	}}).Decode(&createdEndpoint)
	if findErr != nil {
		return nil, GetDbError(findErr)
	}
	return &createdEndpoint, nil
}

func GetEndpoints(mockServiceID string) ([]models.Endpoint, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endpointsCol := Database.Collection(ENDPOINT_COLLECTION)
	dbQuery := bson.D{
		{Key: "mockServiceId", Value: mockServiceID},
	}

	cursor, errFind := endpointsCol.Find(ctx, dbQuery)
	if errFind != nil {
		return []models.Endpoint{}, GetDbError(errFind)
	}
	result := make([]models.Endpoint, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, GetDbError(err)
	}
	return result, nil
}

func BatchLoadEndpoints(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	results := make([]*dataloader.Result, len(keys))
	eCol := Database.Collection(ENDPOINT_COLLECTION)
	dbQuery := bson.D{
		{Key: "id", Value: bson.D{{
			Key: "$in", Value: keys,
		}}},
	}

	cursor, errFind := eCol.Find(ctx, dbQuery)
	if errFind != nil {
		return []*dataloader.Result{{
			Data:  nil,
			Error: errFind,
		}}
	}
	endpoints := make([]models.Endpoint, 0)
	err := cursor.All(ctx, &endpoints)
	if err != nil {
		return []*dataloader.Result{{
			Data:  nil,
			Error: errFind,
		}}
	}

	for i := 0; i < len(keys); i++ {
		results[i] = &dataloader.Result{}
		for _, s := range endpoints {
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

func GetEndpointByID(ID string) (*models.Endpoint, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var ep models.Endpoint
	eCol := Database.Collection(ENDPOINT_COLLECTION)
	findErr := eCol.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: ID,
	}}).Decode(&ep)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, GetDbError(findErr)
	}
	return &ep, nil
}

func DeleteEndpoint(userID, endpointID string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endpointsCol := Database.Collection(ENDPOINT_COLLECTION)
	result, errDel := endpointsCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: endpointID},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, GetDbError(errDel)
	}
	return result.DeletedCount > 0, nil
}
