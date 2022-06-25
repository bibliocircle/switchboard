package db

import (
	"context"
	"strings"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateEndpoint(userId string, ep *models.CreateEndpointRequestBody) (*models.Endpoint, *common.DetailedError) {
	eId, _ := uuid.NewRandom()
	endpointId := eId.String()
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
		return nil, common.WrapAsDetailedError(findErr)
	}
	return &createdEndpoint, nil
}

func GetEndpoints(mockServiceID string) ([]models.Endpoint, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	endpointsCol := Database.Collection(ENDPOINT_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}
	dbQuery := bson.D{
		{Key: "mockServiceId", Value: mockServiceID},
	}

	cursor, errFind := endpointsCol.Find(ctx, dbQuery, findOpts)
	if errFind != nil {
		return []models.Endpoint{}, common.WrapAsDetailedError(errFind)
	}
	result := make([]models.Endpoint, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return result, nil
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
		return false, common.WrapAsDetailedError(errDel)
	}
	return result.DeletedCount > 0, nil
}
