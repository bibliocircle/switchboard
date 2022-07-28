package db

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddMockServiceToWorkspace(userID, workspaceID, mockServiceID string) *common.DetailedError {
	endpoints, errEp := GetEndpoints(mockServiceID)
	if errEp != nil {
		return errEp
	}

	mockService, errMs := GetMockServiceByID(mockServiceID)
	if errMs != nil {
		return errMs
	}

	endpointConfigs := make([]models.EndpointConfig, 0)

	for _, ep := range endpoints {
		sc := make([]models.ScenarioConfig, 0)
		scenarios, errSc := GetScenarios(ep.ID)
		if errSc != nil {
			return common.WrapAsDetailedError(errSc)
		}

		for _, s := range scenarios {
			sc = append(sc, models.ScenarioConfig{
				ScenarioID: s.ID,
				IsActive:   s.IsDefault,
			})
		}

		endpointConfigs = append(endpointConfigs, models.EndpointConfig{
			EndpointID:      ep.ID,
			ScenarioConfigs: sc,
			ResponseDelay:   ep.ResponseDelay,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection(WORKSPACE_SETTINGS_COLLECTION)

	_, insertErr := wssCol.InsertOne(ctx, models.WorkspaceSetting{
		WorkspaceID:     workspaceID,
		MockServiceID:   mockServiceID,
		Config:          mockService.Config,
		EndpointConfigs: endpointConfigs,
	})
	if insertErr != nil {
		return GetDbError(insertErr)
	}

	return nil
}

func GetWorkspaceSettings(workspaceID string) (*[]models.WorkspaceSetting, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection(WORKSPACE_SETTINGS_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": 1,
		},
	}

	cursor, errFind := wssCol.Find(ctx, bson.D{
		{Key: "workspaceId", Value: workspaceID},
	}, findOpts)
	if errFind != nil {
		return nil, common.WrapAsDetailedError(errFind)
	}
	result := make([]models.WorkspaceSetting, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return &result, nil
}

func GetWorkspaceSetting(workspaceID, mockServiceID string) (*models.WorkspaceSetting, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection(WORKSPACE_SETTINGS_COLLECTION)

	var wss models.WorkspaceSetting
	err := wssCol.FindOne(ctx, bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "mockServiceId", Value: mockServiceID},
	}).Decode(&wss)

	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return &wss, nil
}

func ActivateWsMockServiceScenario(workspaceID, mockServiceID, endpointID, scenarioID string) (bool, *common.DetailedError) {
	wss, errWss := GetWorkspaceSetting(workspaceID, mockServiceID)
	if errWss != nil {
		return false, errWss
	}

	newEndpointConfigs := make([]models.EndpointConfig, len(wss.EndpointConfigs))

	for endpointConfigIndex, e := range wss.EndpointConfigs {
		if e.EndpointID == endpointID {
			newScenarioConfigs := make([]models.ScenarioConfig, len(e.ScenarioConfigs))
			for scenarioConfigIndex, sc := range e.ScenarioConfigs {
				if sc.ScenarioID == scenarioID {
					sc.IsActive = true
				} else {
					sc.IsActive = false
				}
				newScenarioConfigs[scenarioConfigIndex] = sc
			}
			newEndpointConfigs[endpointConfigIndex] = models.EndpointConfig{
				EndpointID:      e.EndpointID,
				ScenarioConfigs: newScenarioConfigs,
				ResponseDelay:   e.ResponseDelay,
			}
		} else {
			newEndpointConfigs[endpointConfigIndex] = models.EndpointConfig{
				EndpointID:      e.EndpointID,
				ScenarioConfigs: e.ScenarioConfigs,
				ResponseDelay:   e.ResponseDelay,
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection("workspace_settings")

	result, err := wssCol.UpdateOne(ctx, bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "mockServiceId", Value: mockServiceID},
	}, bson.D{{
		Key: "$set",
		Value: models.WorkspaceSetting{
			WorkspaceID:     wss.WorkspaceID,
			MockServiceID:   wss.MockServiceID,
			Config:          wss.Config,
			EndpointConfigs: newEndpointConfigs,
		},
	}})

	if err != nil {
		return false, common.WrapAsDetailedError(err)
	}

	if result.ModifiedCount == 0 {
		return false, nil
	}

	return true, nil
}
