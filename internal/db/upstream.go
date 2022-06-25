package db

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUpstream(userID string, upstream *models.CreateUpstreamRequestBody) (*models.Upstream, *common.DetailedError) {
	eId, _ := uuid.NewRandom()
	upstreamId := eId.String()
	currentTime := time.Now()
	newUpstream := &models.Upstream{
		ID:            upstreamId,
		MockServiceId: upstream.MockServiceId,
		Name:          upstream.Name,
		URL:           upstream.URL,
		CreatedBy:     userID,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	upstreamsCollection := Database.Collection(UPSTREAMS_COLLECTION)
	_, insertErr := upstreamsCollection.InsertOne(ctx, newUpstream)
	if insertErr != nil {
		return nil, GetDbError(insertErr)
	}
	var createdUpstream models.Upstream
	findErr := upstreamsCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: upstreamId,
	}}).Decode(&createdUpstream)
	if findErr != nil {
		return nil, common.WrapAsDetailedError(findErr)
	}
	return &createdUpstream, nil
}

func GetUpstreams(mockServiceID string) ([]models.Upstream, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	upstreamsCol := Database.Collection(UPSTREAMS_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": -1,
		},
	}
	dbQuery := bson.D{
		{Key: "mockServiceId", Value: mockServiceID},
	}

	cursor, errFind := upstreamsCol.Find(ctx, dbQuery, findOpts)
	if errFind != nil {
		return []models.Upstream{}, common.WrapAsDetailedError(errFind)
	}
	result := make([]models.Upstream, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return result, nil
}

func DeleteUpstream(userID, upstreamID string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	upstreamsCol := Database.Collection(UPSTREAMS_COLLECTION)
	result, errDel := upstreamsCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: upstreamID},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, common.WrapAsDetailedError(errDel)
	}
	return result.DeletedCount > 0, nil
}
