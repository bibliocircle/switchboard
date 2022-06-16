package scenario

import (
	"context"
	"switchboard/internal/common/db"
	"switchboard/internal/common/err_utils"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type ScenarioConfig struct {
	ID                   string            `json:"id" bson:"id,omitempty"`
	StatusCode           int32             `json:"statusCode" bson:"statusCode,omitempty"`
	ResponseBodyTemplate string            `json:"responseBodyTemplate" bson:"responseBodyTemplate,omitempty"`
	ResponseHeaders      map[string]string `json:"responseHeaders" bson:"responseHeaders,omitempty"`
}

type Scenario struct {
	ID         string         `json:"id" bson:"id,omitempty"`
	EndpointId string         `json:"endpointId" bson:"endpointId,omitempty"`
	Type       string         `json:"type" bson:"type,omitempty"`
	Config     ScenarioConfig `json:"config" bson:"config,omitempty"`
	CreatedBy  string         `json:"createdBy" bson:"createdAt,omitempty"`
	CreatedAt  time.Time      `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt  time.Time      `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func CreateScenario(userId string, scenario *Scenario) (*Scenario, *err_utils.DetailedError) {
	eId, _ := uuid.NewRandom()
	scenarioId := eId.String()
	currentTime := time.Now()
	newScenario := &Scenario{
		ID:         scenarioId,
		EndpointId: scenario.EndpointId,
		Type:       scenario.Type,
		Config:     scenario.Config,
		CreatedBy:  userId,
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCollection := db.Database.Collection(db.SCENARIOS_COLLECTION)
	_, insertErr := scenariosCollection.InsertOne(ctx, newScenario)
	if insertErr != nil {
		return nil, db.GetDbError(insertErr)
	}
	var createdScenario Scenario
	findErr := scenariosCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: scenarioId,
	}}).Decode(&createdScenario)
	if findErr != nil {
		return nil, err_utils.WrapAsDetailedError(findErr)
	}
	return &createdScenario, nil
}
