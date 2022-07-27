package db

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func CountScenarios(endpointID string) (int64, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := Database.Collection(SCENARIOS_COLLECTION)

	count, errCount := scenariosCol.CountDocuments(ctx, bson.D{
		{Key: "endpointId", Value: endpointID},
	})
	if errCount != nil {
		return 0, common.WrapAsDetailedError(errCount)
	}
	return count, nil
}

func CreateScenario(userId string, sc *models.CreateScenarioRequestBody) (*models.Scenario, *common.DetailedError) {
	eId, _ := uuid.NewRandom()
	scenarioId := eId.String()
	currentTime := time.Now()
	isDefaultScenario := false
	count, err := CountScenarios(sc.EndpointId)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	isDefaultScenario = count == 0
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
	case common.HTTP_SCENARIO_TYPE:
		newScenario.HTTPResponseScenarioConfig = &models.HTTPResponseScenarioConfig{
			StatusCode:              sc.HTTPResponseScenarioConfig.StatusCode,
			ResponseBodyTemplate:    sc.HTTPResponseScenarioConfig.ResponseBodyTemplate,
			ResponseHeadersTemplate: sc.HTTPResponseScenarioConfig.ResponseHeadersTemplate,
		}
	case common.PROXY_SCENARIO_TYPE:
		newScenario.ProxyScenarioConfig = &models.ProxyScenarioConfig{
			Name:          sc.ProxyScenarioConfig.Name,
			UpstreamID:    sc.ProxyScenarioConfig.UpstreamID,
			InjectHeaders: sc.ProxyScenarioConfig.InjectHeaders,
		}
	case common.NETWORK_SCENARIO_TYPE:
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
		return nil, common.WrapAsDetailedError(findErr)
	}
	return &createdScenario, nil
}

func GetScenarios(endpointID string) ([]models.Scenario, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := Database.Collection(SCENARIOS_COLLECTION)
	dbQuery := bson.D{
		{Key: "endpointId", Value: endpointID},
	}

	cursor, errFind := scenariosCol.Find(ctx, dbQuery)
	if errFind != nil {
		return nil, common.WrapAsDetailedError(errFind)
	}
	result := make([]models.Scenario, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return result, nil
}

func GetScenarioByID(scenarioID string) (*models.Scenario, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := Database.Collection(SCENARIOS_COLLECTION)

	var sc models.Scenario
	err := scenariosCol.FindOne(ctx, bson.D{
		{Key: "id", Value: scenarioID},
	}).Decode(&sc)

	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return &sc, nil
}
