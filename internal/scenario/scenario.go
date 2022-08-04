package scenario

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/util"
	"time"

	"github.com/graph-gophers/dataloader"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CountScenarios(endpointID string) (int64, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := db.Database.Collection(db.SCENARIOS_COLLECTION)

	count, errCount := scenariosCol.CountDocuments(ctx, bson.D{
		{Key: "endpointId", Value: endpointID},
	})
	if errCount != nil {
		return 0, db.GetDbError(errCount)
	}
	return count, nil
}

func CreateScenario(userId string, sc *CreateScenarioRequestBody) (*Scenario, *common.DetailedError) {
	scenarioId := util.UUIDv4()
	currentTime := time.Now()
	isDefaultScenario := false
	count, err := CountScenarios(sc.EndpointId)
	if err != nil {
		return nil, db.GetDbError(err)
	}
	isDefaultScenario = count == 0
	newScenario := &Scenario{
		ID:         scenarioId,
		EndpointId: sc.EndpointId,
		Type:       sc.Type,
		IsDefault:  isDefaultScenario,
		CreatedBy:  userId,
		CreatedAt:  &currentTime,
		UpdatedAt:  &currentTime,
	}
	switch sc.Type {
	case common.HTTP_SCENARIO_TYPE:
		newScenario.HTTPResponseScenarioConfig = &HTTPResponseScenarioConfig{
			StatusCode:              sc.HTTPResponseScenarioConfig.StatusCode,
			ResponseBodyTemplate:    sc.HTTPResponseScenarioConfig.ResponseBodyTemplate,
			ResponseHeadersTemplate: sc.HTTPResponseScenarioConfig.ResponseHeadersTemplate,
		}
	case common.PROXY_SCENARIO_TYPE:
		newScenario.ProxyScenarioConfig = &ProxyScenarioConfig{
			UpstreamID:    sc.ProxyScenarioConfig.UpstreamID,
			InjectHeaders: sc.ProxyScenarioConfig.InjectHeaders,
		}
	case common.NETWORK_SCENARIO_TYPE:
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
		if findErr == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, db.GetDbError(findErr)
	}
	return &createdScenario, nil
}

func GetScenarios(endpointID string) ([]Scenario, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := db.Database.Collection(db.SCENARIOS_COLLECTION)
	dbQuery := bson.D{
		{Key: "endpointId", Value: endpointID},
	}

	cursor, errFind := scenariosCol.Find(ctx, dbQuery)
	if errFind != nil {
		return nil, db.GetDbError(errFind)
	}
	result := make([]Scenario, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, db.GetDbError(err)
	}
	return result, nil
}

func BatchLoadScenarios(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	results := make([]*dataloader.Result, len(keys))
	scenariosCol := db.Database.Collection(db.SCENARIOS_COLLECTION)
	dbQuery := bson.D{
		{Key: "id", Value: bson.D{{
			Key: "$in", Value: keys,
		}}},
	}

	cursor, errFind := scenariosCol.Find(ctx, dbQuery)
	if errFind != nil {
		return []*dataloader.Result{{
			Data:  nil,
			Error: errFind,
		}}
	}
	scenarios := make([]Scenario, 0)
	err := cursor.All(ctx, &scenarios)
	if err != nil {
		return []*dataloader.Result{{
			Data:  nil,
			Error: errFind,
		}}
	}

	for i := 0; i < len(keys); i++ {
		results[i] = &dataloader.Result{}
		for _, s := range scenarios {
			if s.ID == keys[i].String() {
				results[i] = &dataloader.Result{
					Data:  &s,
					Error: nil,
				}
				break
			}
		}
	}

	return results
}

func GetScenarioByID(scenarioID string) (*Scenario, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	scenariosCol := db.Database.Collection(db.SCENARIOS_COLLECTION)

	var sc Scenario
	err := scenariosCol.FindOne(ctx, bson.D{
		{Key: "id", Value: scenarioID},
	}).Decode(&sc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, db.GetDbError(err)
	}
	return &sc, nil
}
