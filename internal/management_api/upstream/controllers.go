package upstream

import (
	"context"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/db"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Upstream struct {
	ID            string    `json:"id" bson:"id,omitempty"`
	MockServiceId string    `json:"mockServiceId" bson:"mockServiceId,omitempty"`
	Name          string    `json:"name" bson:"name,omitempty"`
	URL           string    `json:"url" bson:"url,omitempty"`
	CreatedBy     string    `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func CreateUpstream(userID string, upstream *CreateUpstreamRequestBody) (*Upstream, *err_utils.DetailedError) {
	eId, _ := uuid.NewRandom()
	upstreamId := eId.String()
	currentTime := time.Now()
	newUpstream := &Upstream{
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
	upstreamsCollection := db.Database.Collection(db.UPSTREAMS_COLLECTION)
	_, insertErr := upstreamsCollection.InsertOne(ctx, newUpstream)
	if insertErr != nil {
		return nil, db.GetDbError(insertErr)
	}
	var createdUpstream Upstream
	findErr := upstreamsCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: upstreamId,
	}}).Decode(&createdUpstream)
	if findErr != nil {
		return nil, err_utils.WrapAsDetailedError(findErr)
	}
	return &createdUpstream, nil
}

func GetUpstreams(mockServiceID string) ([]Upstream, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	upstreamsCol := db.Database.Collection(db.UPSTREAMS_COLLECTION)
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
		return []Upstream{}, err_utils.WrapAsDetailedError(errFind)
	}
	result := make([]Upstream, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	return result, nil
}

func DeleteUpstream(userID, upstreamID string) (bool, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	upstreamsCol := db.Database.Collection(db.UPSTREAMS_COLLECTION)
	result, errDel := upstreamsCol.DeleteOne(ctx, bson.D{
		{Key: "id", Value: upstreamID},
		{Key: "createdBy", Value: userID},
	})
	if errDel != nil {
		return false, err_utils.WrapAsDetailedError(errDel)
	}
	return result.DeletedCount > 0, nil
}
