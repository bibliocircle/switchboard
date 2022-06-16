package upstream

import (
	"context"
	"switchboard/internal/common/db"
	"switchboard/internal/common/err_utils"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

func CreateUpstream(userId string, upstream *CreateUpstreamRequestBody) (*Upstream, *err_utils.DetailedError) {
	eId, _ := uuid.NewRandom()
	upstreamId := eId.String()
	currentTime := time.Now()
	newUpstream := &Upstream{
		ID:            upstreamId,
		MockServiceId: upstream.MockServiceId,
		Name:          upstream.Name,
		URL:           upstream.URL,
		CreatedBy:     userId,
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
