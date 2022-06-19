package workspace_settings

import (
	"context"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/db"
	"switchboard/internal/management_api/endpoint"
	"switchboard/internal/management_api/scenario"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScenarioConfig struct {
	ScenarioID string `json:"scenarioId" bson:"id,omitempty"`
	IsActive   bool   `json:"isActive" bson:"isActive"`
}

type WorkspaceSetting struct {
	WorkspaceID   string           `json:"workspaceId" bson:"workspaceId,omitempty"`
	MockServiceID string           `json:"mockServiceId" bson:"mockServiceId,omitempty"`
	EndpointID    string           `json:"endpointId" bson:"endpointId,omitempty"`
	Scenarios     []ScenarioConfig `json:"scenarios" bson:"scenarios"`
	ResponseDelay uint16           `json:"responseDelay" bson:"responseDelay,omitempty"`
	CreatedBy     string           `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt     time.Time        `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time        `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func CreateWorkspaceSetting(userID, workspaceID, mockServiceID, endpointID string, scenarios []ScenarioConfig) *err_utils.DetailedError {
	currentTime := time.Now()
	newWss := &WorkspaceSetting{
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
	wssCol := db.Database.Collection(db.WORKSPACE_SETTINGS_COLLECTION)
	_, insertErr := wssCol.InsertOne(ctx, newWss)
	if insertErr != nil {
		return db.GetDbError(insertErr)
	}
	return nil
}

func AddMockServiceToWorkspace(userID, workspaceID, mockServiceID string) *err_utils.DetailedError {
	endpoints, errEp := endpoint.GetEndpoints(mockServiceID)
	if errEp != nil {
		return err_utils.WrapAsDetailedError(errEp)
	}

	// TODO: This can be improved to use goroutines and channels
	for _, ep := range endpoints {
		sc := make([]ScenarioConfig, 0)
		scenarios, errSc := scenario.GetScenarios(ep.ID)
		if errSc != nil {
			return err_utils.WrapAsDetailedError(errSc)
		}

		for _, s := range scenarios {
			sc = append(sc, ScenarioConfig{
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

func GetWorkspaceSettings(workspaceID string) ([]WorkspaceSetting, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := db.Database.Collection(db.WORKSPACE_SETTINGS_COLLECTION)
	findOpts := &options.FindOptions{
		Sort: &map[string]int64{
			"createdAt": 1,
		},
	}

	cursor, errFind := wssCol.Find(ctx, bson.D{
		{Key: "workspaceId", Value: workspaceID},
	}, findOpts)
	if errFind != nil {
		return nil, err_utils.WrapAsDetailedError(errFind)
	}
	result := make([]WorkspaceSetting, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	return result, nil
}

func ActivateMockServiceScenario(workspaceID, mockServiceID, endpointID, scenarioID string) (bool, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := db.Database.Collection("workspace_settings")
	q := bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "mockServiceId", Value: mockServiceID},
		{Key: "endpointId", Value: endpointID},
	}
	var wss WorkspaceSetting
	findErr := wssCol.FindOne(ctx, q).Decode(&wss)
	if findErr != nil {
		return false, err_utils.WrapAsDetailedError(findErr)
	}

	scenarios := make([]ScenarioConfig, 0)
	for _, s := range wss.Scenarios {
		if s.ScenarioID == scenarioID {
			scenarios = append(scenarios, ScenarioConfig{
				ScenarioID: s.ScenarioID,
				IsActive:   true,
			})
		} else {
			scenarios = append(scenarios, ScenarioConfig{
				ScenarioID: s.ScenarioID,
				IsActive:   false,
			})
		}
	}

	result, err := wssCol.UpdateOne(ctx, q, bson.D{{
		Key: "$set",
		Value: &WorkspaceSetting{
			Scenarios: scenarios,
			UpdatedAt: time.Now(),
		},
	}})
	if err != nil {
		return false, err_utils.WrapAsDetailedError(err)
	}
	if result.ModifiedCount == 0 {
		return false, nil
	}

	return true, nil
}

func UpdateWsMockServiceConfig(workspaceID, mockServiceID, endpointID string, wssUpdate *UpdateMockServiceConfigRequestBody) (*WorkspaceSetting, *err_utils.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := db.Database.Collection("workspace_settings")
	q := bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "mockServiceId", Value: mockServiceID},
		{Key: "endpointId", Value: endpointID},
	}
	result, err := wssCol.UpdateOne(ctx, q, bson.D{{
		Key: "$set",
		Value: &WorkspaceSetting{
			ResponseDelay: wssUpdate.ResponseDelay,
			UpdatedAt:     time.Now(),
		},
	}})
	if err != nil {
		return nil, err_utils.WrapAsDetailedError(err)
	}
	if result.ModifiedCount == 0 {
		return nil, nil
	}
	var updatedWss WorkspaceSetting
	findErr := wssCol.FindOne(ctx, q).Decode(&updatedWss)
	if findErr != nil {
		return nil, err_utils.WrapAsDetailedError(findErr)
	}
	return &updatedWss, nil
}
