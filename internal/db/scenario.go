package db

import (
	"context"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func CountScenarios(endpointID string) (int64, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := Database.Collection(SCENARIOS_COLLECTION)

	count, errCount := scenariosCol.CountDocuments(ctx, bson.D{
		{Key: "endpointId", Value: endpointID},
	})
	if errCount != nil {
		return 0, err_utils.WrapAsDetailedError(errCount)
	}
	return count, nil
}

func CreateScenario(userId string, sc *models.CreateScenarioRequestBody) (*models.Scenario, *err_utils.DetailedError) {
	eId, _ := uuid.NewRandom()
	scenarioId := eId.String()
	currentTime := time.Now()
	isDefaultScenario := false
	count, err := CountScenarios(sc.EndpointId)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	isDefaultScenario = count != 0
	newScenario := &models.Scenario{
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
		newScenario.HTTPResponseScenarioConfig = &models.HTTPResponseScenarioConfig{
			StatusCode:              sc.HTTPResponseScenarioConfig.StatusCode,
			ResponseBodyTemplate:    sc.HTTPResponseScenarioConfig.ResponseBodyTemplate,
			ResponseHeadersTemplate: sc.HTTPResponseScenarioConfig.ResponseHeadersTemplate,
		}
	case constants.PROXY_SCENARIO_TYPE:
		newScenario.ProxyScenarioConfig = &models.ProxyScenarioConfig{
			UpstreamID:    sc.ProxyScenarioConfig.UpstreamID,
			InjectHeaders: sc.ProxyScenarioConfig.InjectHeaders,
		}
	case constants.NETWORK_SCENARIO_TYPE:
		newScenario.NetworkScenarioConfig = &models.NetworkScenarioConfig{
			Type: sc.NetworkScenarioConfig.Type,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCollection := Database.Collection(SCENARIOS_COLLECTION)
	_, insertErr := scenariosCollection.InsertOne(ctx, newScenario)
	if insertErr != nil {
		return nil, GetDbError(insertErr)
	}
	var createdScenario models.Scenario
	findErr := scenariosCollection.FindOne(ctx, bson.D{{
		Key:   "id",
		Value: scenarioId,
	}}).Decode(&createdScenario)
	if findErr != nil {
		return nil, err_utils.WrapAsDetailedError(findErr)
	}
	return &createdScenario, nil
}

func GetScenarios(endpointID string) ([]models.Scenario, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := Database.Collection(SCENARIOS_COLLECTION)
	dbQuery := bson.D{
		{Key: "endpointId", Value: endpointID},
	}

	cursor, errFind := scenariosCol.Find(ctx, dbQuery)
	if errFind != nil {
		return nil, err_utils.WrapAsDetailedError(errFind)
	}
	result := make([]models.Scenario, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	return result, nil
}
