package scenario

import (
	"context"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/db"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type HTTPResponseScenarioConfig struct {
	StatusCode              uint16 `json:"statusCode" bson:"statusCode,omitempty"`
	ResponseBodyTemplate    string `json:"responseBodyTemplate" bson:"responseBodyTemplate,omitempty"`
	ResponseHeadersTemplate string `json:"responseHeadersTemplate" bson:"responseHeadersTemplate,omitempty"`
}

type ProxyScenarioConfig struct {
	UpstreamID    string            `json:"upstreamID" bson:"upstreamID,omitempty"`
	InjectHeaders map[string]string `json:"injectHeaders" bson:"injectHeaders,omitempty"`
}

type NetworkScenarioConfig struct {
	Type string `json:"type" bson:"type,omitempty"`
}

type Scenario struct {
	ID                         string                      `json:"id" bson:"id,omitempty"`
	EndpointId                 string                      `json:"endpointId" bson:"endpointId,omitempty"`
	Type                       string                      `json:"type" bson:"type,omitempty"`
	IsDefault                  bool                        `json:"isDefault" bson:"isDefault"`
	HTTPResponseScenarioConfig *HTTPResponseScenarioConfig `json:"httpResponseScenarioConfig,omitempty" bson:"httpResponseScenarioConfig,omitempty"`
	ProxyScenarioConfig        *ProxyScenarioConfig        `json:"proxyScenarioConfig,omitempty" bson:"proxyScenarioConfig,omitempty"`
	NetworkScenarioConfig      *NetworkScenarioConfig      `json:"networkScenarioConfig,omitempty" bson:"networkScenarioConfig,omitempty"`
	CreatedBy                  string                      `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt                  time.Time                   `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt                  time.Time                   `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func CountScenarios(endpointID string) (int64, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := db.Database.Collection(db.SCENARIOS_COLLECTION)

	count, errCount := scenariosCol.CountDocuments(ctx, bson.D{
		{Key: "endpointId", Value: endpointID},
	})
	if errCount != nil {
		return 0, err_utils.WrapAsDetailedError(errCount)
	}
	return count, nil
}

func CreateScenario(userId string, sc *CreateScenarioRequestBody) (*Scenario, *err_utils.DetailedError) {
	eId, _ := uuid.NewRandom()
	scenarioId := eId.String()
	currentTime := time.Now()
	isDefaultScenario := false
	count, err := CountScenarios(sc.EndpointId)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	isDefaultScenario = count != 0
	newScenario := &Scenario{
		ID:         scenarioId,
		EndpointId: sc.EndpointId,
		Type:       sc.Type,
		IsDefault:  isDefaultScenario,
		CreatedBy:  userId,
		CreatedAt:  currentTime,
		UpdatedAt:  currentTime,
	}
	switch sc.Type {
	case constants.HTTP_SCENARIO_TYPE:
		newScenario.HTTPResponseScenarioConfig = &HTTPResponseScenarioConfig{
			StatusCode:              sc.HTTPResponseScenarioConfig.StatusCode,
			ResponseBodyTemplate:    sc.HTTPResponseScenarioConfig.ResponseBodyTemplate,
			ResponseHeadersTemplate: sc.HTTPResponseScenarioConfig.ResponseHeadersTemplate,
		}
	case constants.PROXY_SCENARIO_TYPE:
		newScenario.ProxyScenarioConfig = &ProxyScenarioConfig{
			UpstreamID:    sc.ProxyScenarioConfig.UpstreamID,
			InjectHeaders: sc.ProxyScenarioConfig.InjectHeaders,
		}
	case constants.NETWORK_SCENARIO_TYPE:
		newScenario.NetworkScenarioConfig = &NetworkScenarioConfig{
			Type: sc.NetworkScenarioConfig.Type,
		}
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

func GetScenarios(endpointID string) ([]Scenario, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := db.Database.Collection(db.SCENARIOS_COLLECTION)
	dbQuery := bson.D{
		{Key: "endpointId", Value: endpointID},
	}

	cursor, errFind := scenariosCol.Find(ctx, dbQuery)
	if errFind != nil {
		return nil, err_utils.WrapAsDetailedError(errFind)
	}
	result := make([]Scenario, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	return result, nil
}
