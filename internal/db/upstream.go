package db

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"switchboard/internal/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUpstream(userID string, upstream *models.CreateUpstreamRequestBody) (*models.Upstream, *common.DetailedError) {
	upstreamId := util.UUIDv4()
	currentTime := time.Now()
	newUpstream := &models.Upstream{
		ID:            upstreamId,
		MockServiceId: upstream.MockServiceId,
		Name:          upstream.Name,
		URL:           upstream.URL,
		CreatedBy:     userID,
		CreatedAt:     &currentTime,
		UpdatedAt:     &currentTime,
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
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, GetDbError(findErr)
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
		return []models.Upstream{}, GetDbError(errFind)
	}
	result := make([]models.Upstream, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, GetDbError(err)
	}
	return result, nil
}

func GetUpstreamByID(ID string) (*models.Upstream, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var upstream models.Upstream
	upstreamsCol := Database.Collection(UPSTREAMS_COLLECTION)
	findErr := upstreamsCol.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: ID,
	}}).Decode(&upstream)
	if findErr != nil {
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, GetDbError(findErr)
	}
	return &upstream, nil
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
		return false, GetDbError(errDel)
	}
	return result.DeletedCount > 0, nil
}
