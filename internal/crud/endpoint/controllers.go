package endpoint

import (
	"context"
	"strings"
	"switchboard/internal/common/db"
	"switchboard/internal/common/err_utils"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Endpoint struct {
	ID            string    `json:"id" bson:"id,omitempty"`
	MockServiceId string    `json:"mockServiceId" bson:"mockServiceId,omitempty"`
	Path          string    `json:"path" bson:"path,omitempty"`
	Method        string    `json:"method" bson:"method,omitempty"`
	Description   string    `json:"description" bson:"description,omitempty"`
	ResponseDelay int64     `json:"responseDelay" bson:"responseDelay,omitempty"`
	CreatedBy     string    `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func CreateEndpoint(userId string, ep *CreateEndpointRequestBody) (*Endpoint, *err_utils.DetailedError) {
	eId, _ := uuid.NewRandom()
	endpointId := eId.String()
	currentTime := time.Now()
	newEndpoint := &Endpoint{
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
	endpointsCollection := db.Database.Collection(db.ENDPOINT_COLLECTION)
	_, insertErr := endpointsCollection.InsertOne(ctx, newEndpoint)
	if insertErr != nil {
		return nil, db.GetDbError(insertErr)
	}
	var createdEndpoint Endpoint
	findErr := endpointsCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: endpointId,
	}}).Decode(&createdEndpoint)
	if findErr != nil {
		return nil, err_utils.WrapAsDetailedError(findErr)
	}
	return &createdEndpoint, nil
}
