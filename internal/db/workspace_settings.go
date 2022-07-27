package db

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateWorkspaceSetting(userID, workspaceID, mockServiceID, endpointID string, scenarios []models.ScenarioConfig) *common.DetailedError {
	currentTime := time.Now()
	newWss := &models.WorkspaceSetting{
		WorkspaceID:   workspaceID,
		MockServiceID: mockServiceID,
		EndpointID:    endpointID,
		Scenarios:     scenarios,
		ResponseDelay: 0,
		CreatedBy:     userID,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection(WORKSPACE_SETTINGS_COLLECTION)
	_, insertErr := wssCol.InsertOne(ctx, newWss)
	if insertErr != nil {
		return GetDbError(insertErr)
	}
	return nil
}

func AddMockServiceToWorkspace(userID, workspaceID, mockServiceID string) *common.DetailedError {
	endpoints, errEp := GetEndpoints(mockServiceID)
	if errEp != nil {
		return common.WrapAsDetailedError(errEp)
	}

	// TODO: This can be improved to use goroutines and channels
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

		errCreateWss := CreateWorkspaceSetting(userID, workspaceID, mockServiceID, ep.ID, sc)
		if errCreateWss != nil {
			// TODO: delete any partially inserted settings
			return errCreateWss
		}
	}
	return nil
}

func GetWorkspaceSettings(workspaceID string) ([]models.WorkspaceSetting, *common.DetailedError) {
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
	return result, nil
}

func GetWorkspaceMockServices(workspaceID string) ([]models.MockService, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection(WORKSPACE_SETTINGS_COLLECTION)

	mockServiceIDs, errMSIDs := wssCol.Distinct(ctx, "mockServiceId", bson.D{
		{Key: "workspaceId", Value: workspaceID},
	})
	if errMSIDs != nil {
		return nil, common.WrapAsDetailedError(errMSIDs)
	}

	msFindOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": 1,
		},
	}
	msCollection := Database.Collection(MOCKSERVICES_COLLECTION)
	cursor, errFind := msCollection.Find(ctx, bson.D{
		{Key: "id", Value: bson.D{{
			Key:   "$in",
			Value: mockServiceIDs},
		}},
	}, msFindOpts)
	if errFind != nil {
		return nil, common.WrapAsDetailedError(errFind)
	}
	mockServices := make([]models.MockService, 0)
	err := cursor.All(ctx, &mockServices)
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return mockServices, nil
}

func GetWorkspaceSetting(workspaceID, endpointID string) (*models.WorkspaceSetting, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection(WORKSPACE_SETTINGS_COLLECTION)

	var wss models.WorkspaceSetting
	err := wssCol.FindOne(ctx, bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "endpointId", Value: endpointID},
	}).Decode(&wss)

	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	return &wss, nil
}

func ActivateMockServiceScenario(workspaceID, mockServiceID, endpointID, scenarioID string) (bool, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection("workspace_settings")
	q := bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "mockServiceId", Value: mockServiceID},
		{Key: "endpointId", Value: endpointID},
	}
	var wss models.WorkspaceSetting
	findErr := wssCol.FindOne(ctx, q).Decode(&wss)
	if findErr != nil {
		return false, common.WrapAsDetailedError(findErr)
	}

	scenarios := make([]models.ScenarioConfig, 0)
	for _, s := range wss.Scenarios {
		if s.ScenarioID == scenarioID {
			scenarios = append(scenarios, models.ScenarioConfig{
				ScenarioID: s.ScenarioID,
				IsActive:   true,
			})
		} else {
			scenarios = append(scenarios, models.ScenarioConfig{
				ScenarioID: s.ScenarioID,
				IsActive:   false,
			})
		}
	}

	result, err := wssCol.UpdateOne(ctx, q, bson.D{{
		Key: "$set",
		Value: &models.WorkspaceSetting{
			Scenarios: scenarios,
			UpdatedAt: time.Now(),
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

func UpdateWsMockServiceConfig(workspaceID, mockServiceID, endpointID string, wssUpdate *models.UpdateMockServiceConfigRequestBody) (*models.WorkspaceSetting, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := Database.Collection("workspace_settings")
	q := bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "mockServiceId", Value: mockServiceID},
		{Key: "endpointId", Value: endpointID},
	}
	result, err := wssCol.UpdateOne(ctx, q, bson.D{{
		Key: "$set",
		Value: &models.WorkspaceSetting{
			ResponseDelay: wssUpdate.ResponseDelay,
			UpdatedAt:     time.Now(),
		},
	}})
	if err != nil {
		return nil, common.WrapAsDetailedError(err)
	}
	if result.ModifiedCount == 0 {
		return nil, nil
	}
	var updatedWss models.WorkspaceSetting
	findErr := wssCol.FindOne(ctx, q).Decode(&updatedWss)
	if findErr != nil {
		return nil, common.WrapAsDetailedError(findErr)
	}
	return &updatedWss, nil
}
