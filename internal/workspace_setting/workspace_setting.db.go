package workspace_setting

import (
	"context"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/endpoint"
	"switchboard/internal/mockservice"
	"switchboard/internal/scenario"
	"switchboard/internal/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddMockServiceToWorkspace(userID, workspaceID, mockServiceID string) *common.DetailedError {
	endpoints, errEp := endpoint.GetEndpoints(mockServiceID)
	if errEp != nil {
		return errEp
	}

	ms, errMs := mockservice.GetMockServiceByID(mockServiceID)
	if errMs != nil {
		return errMs
	}

	endpointConfigs := make([]EndpointConfig, 0)

	for _, ep := range endpoints {
		sc := make([]ScenarioConfig, 0)
		scenarios, errSc := scenario.GetScenarios(ep.ID)
		if errSc != nil {
			return db.GetDbError(errSc)
		}

		for _, s := range scenarios {
			sc = append(sc, ScenarioConfig{
				ID:         util.UUIDv4(),
				ScenarioID: s.ID,
				IsActive:   s.IsDefault,
			})
		}

		endpointConfigs = append(endpointConfigs, EndpointConfig{
			ID:              util.UUIDv4(),
			EndpointID:      ep.ID,
			ScenarioConfigs: sc,
			ResponseDelay:   ep.ResponseDelay,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := db.Database.Collection(db.WORKSPACE_SETTINGS_COLLECTION)

	_, insertErr := wssCol.InsertOne(ctx, WorkspaceSetting{
		ID:              util.UUIDv4(),
		WorkspaceID:     workspaceID,
		MockServiceID:   mockServiceID,
		Config:          ms.Config,
		EndpointConfigs: endpointConfigs,
	})
	if insertErr != nil {
		return db.GetDbError(insertErr)
	}

	return nil
}

func GetWorkspaceSettings(workspaceID string) (*[]*WorkspaceSetting, *common.DetailedError) {
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
		return nil, db.GetDbError(errFind)
	}
	result := make([]*WorkspaceSetting, 0)
	err := cursor.All(ctx, &result)
	if err != nil {
		return nil, db.GetDbError(err)
	}
	return &result, nil
}

func GetWorkspaceSetting(workspaceID, mockServiceID string) (*WorkspaceSetting, *common.DetailedError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := db.Database.Collection(db.WORKSPACE_SETTINGS_COLLECTION)

	var wss WorkspaceSetting
	err := wssCol.FindOne(ctx, bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "mockServiceId", Value: mockServiceID},
	}).Decode(&wss)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, db.GetDbError(err)
	}
	return &wss, nil
}

func ActivateWsMockServiceScenario(workspaceID, mockServiceID, endpointID, scenarioID string) (bool, *common.DetailedError) {
	wss, errWss := GetWorkspaceSetting(workspaceID, mockServiceID)
	if errWss != nil {
		return false, errWss
	}

	newEndpointConfigs := make([]EndpointConfig, len(wss.EndpointConfigs))

	for endpointConfigIndex, e := range wss.EndpointConfigs {
		if e.EndpointID == endpointID {
			newScenarioConfigs := make([]ScenarioConfig, len(e.ScenarioConfigs))
			for scenarioConfigIndex, sc := range e.ScenarioConfigs {
				if sc.ScenarioID == scenarioID {
					sc.IsActive = true
				} else {
					sc.IsActive = false
				}
				newScenarioConfigs[scenarioConfigIndex] = sc
			}
			newEndpointConfigs[endpointConfigIndex] = EndpointConfig{
				ID:              e.ID,
				EndpointID:      e.EndpointID,
				ScenarioConfigs: newScenarioConfigs,
				ResponseDelay:   e.ResponseDelay,
			}
		} else {
			newEndpointConfigs[endpointConfigIndex] = EndpointConfig{
				ID:              e.ID,
				EndpointID:      e.EndpointID,
				ScenarioConfigs: e.ScenarioConfigs,
				ResponseDelay:   e.ResponseDelay,
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wssCol := db.Database.Collection("workspace_settings")

	_, err := wssCol.UpdateOne(ctx, bson.D{
		{Key: "workspaceId", Value: workspaceID},
		{Key: "mockServiceId", Value: mockServiceID},
	}, bson.D{{
		Key: "$set",
		Value: WorkspaceSetting{
			ID:              wss.ID,
			WorkspaceID:     wss.WorkspaceID,
			MockServiceID:   wss.MockServiceID,
			Config:          wss.Config,
			EndpointConfigs: newEndpointConfigs,
		},
	}})

	if err != nil {
		return false, db.GetDbError(err)
	}

	return true, nil
}
