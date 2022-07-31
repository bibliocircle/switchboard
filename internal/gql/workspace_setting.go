package gql

import (
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var ScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ScenarioConfig",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"scenario": &graphql.Field{
			Type: ScenarioGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				scenarioID := p.Source.(models.ScenarioConfig).ScenarioID
				loaders := p.Context.Value(db.LoadersCtxKey).(*db.Loaders)
				thunk := loaders.Scenarios.Load(p.Context, dataloader.StringKey(scenarioID))
				return func() (interface{}, error) {
					scenario, err := thunk()
					if err != nil {
						logrus.Errorln(err)
						return nil, NewGqlError(common.ErrorGeneric, "could not retrieve scenario")
					}
					return scenario.(*models.Scenario), nil
				}, nil
			},
		},
		"isActive": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var EndpointConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "EndpointConfig",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"endpoint": &graphql.Field{
			Type:    EndpointGqlType,
			Resolve: GetEndpointResolver,
		},
		"scenarioConfigs": &graphql.Field{
			Type: graphql.NewList(ScenarioConfigGqlType),
		},
		"responseDelay": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var WorkspaceSettingGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "WorkspaceSetting",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"workspace": &graphql.Field{
			Type: WorkspaceGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				workspaceID := p.Source.(*models.WorkspaceSetting).WorkspaceID
				currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
				ws, err := db.GetUserWorkspaceByID(currentUser.ID, workspaceID)
				if err != nil {
					logrus.Errorln(err)
					return nil, NewGqlError(common.ErrorGeneric, "could not retrieve workspace")
				}
				return ws, nil
			},
		},
		"mockService": &graphql.Field{
			Type: MockServiceGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mockServiceID := p.Source.(*models.WorkspaceSetting).MockServiceID
				ms, err := db.GetMockServiceByID(mockServiceID)
				if err != nil {
					logrus.Errorln(err)
					return nil, NewGqlError(common.ErrorGeneric, "could not retrieve mock service")
				}
				return ms, nil
			},
		},
		"config": &graphql.Field{
			Type: GlobalMockServiceConfigGqlType,
		},
		"endpointConfigs": &graphql.Field{
			Type: graphql.NewList(EndpointConfigGqlType),
		},
	},
})

func GetEndpointResolver(p graphql.ResolveParams) (interface{}, error) {
	endpointID := p.Source.(models.EndpointConfig).EndpointID
	loaders := p.Context.Value(db.LoadersCtxKey).(*db.Loaders)
	thunk := loaders.Endpoints.Load(p.Context, dataloader.StringKey(endpointID))
	return func() (interface{}, error) {
		endpoint, err := thunk()
		if err != nil {
			logrus.Errorln(err)
			return nil, NewGqlError(common.ErrorGeneric, "could not retrieve endpoint")
		}
		return endpoint.(*models.Endpoint), nil
	}, nil
}

func GetWorkspaceSettingsResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Source.(models.Workspace).ID
	wss, err := db.GetWorkspaceSettings(workspaceID)
	if err != nil {
		logrus.Errorf("could not retrieve workspace settings for workspace ID %s : %s\n", workspaceID, err)
		return make([]models.WorkspaceSetting, 0), NewGqlError(common.ErrorGeneric, "could not retrieve workspace settings")
	}

	return wss, nil
}

func GetWorkspaceMockServicesResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Source.(models.Workspace).ID
	wss, errWs := db.GetWorkspaceSettings(workspaceID)
	if errWs != nil {
		if errWs.ErrorCode == common.ErrorNotFound {
			return make([]*models.MockService, 0), NewGqlError(common.ErrorNotFound, "workspace settings not found!")
		}
		return make([]*models.MockService, 0), errWs
	}
	mockServiceIds := make([]string, 0)
	for _, ws := range *wss {
		mockServiceIds = append(mockServiceIds, ws.MockServiceID)
	}

	ms, errMs := db.GetMockServicesByIds(mockServiceIds)
	if errMs != nil {
		return nil, errMs
	}

	return ms, nil
}

func GetWorkspaceSettingResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	mockServiceID := p.Args["mockServiceId"].(string)
	wss, err := db.GetWorkspaceSetting(workspaceID, mockServiceID)
	if err == nil {
		return wss, nil
	}
	if err.ErrorCode == common.ErrorNotFound {
		return nil, NewGqlError(common.ErrorNotFound, "workspace not found or mock service not found in the workspace")
	}
	return nil, NewGqlError(common.ErrorGeneric, "could not retrieve workspace settings")
}

func ActivateMockServiceScenarioResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	mockServiceID := p.Args["mockServiceId"].(string)
	endpointID := p.Args["endpointId"].(string)
	scenarioID := p.Args["scenarioId"].(string)

	ok, errAct := db.ActivateWsMockServiceScenario(
		workspaceID,
		mockServiceID,
		endpointID,
		scenarioID,
	)
	if errAct != nil {
		return nil, NewGqlError(common.ErrorGeneric, "could not activate mock service scenario")
	}
	if ok {
		wss, err := db.GetWorkspaceSetting(workspaceID, mockServiceID)
		if err != nil {
			return nil, err
		}
		return wss, nil
	}
	return nil, nil
}

func AddMockServiceToWorkspaceResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	mockServiceID := p.Args["mockServiceId"].(string)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)

	isWsOwner, errPerm := db.IsWorkspaceOwner(currentUser.ID, workspaceID)
	if errPerm != nil {
		logrus.Errorln("could not check workspace permissions", errPerm)
		return nil, NewGqlError(common.ErrorGeneric, "could not check workspace permissions")
	}

	if !isWsOwner {
		return nil, NewGqlError(common.ErrorForbidden, "workspace doesn't exist or the current user doesn't have permissions to access the workspace")
	}

	err := db.AddMockServiceToWorkspace(currentUser.ID, workspaceID, mockServiceID)
	if err == nil {
		ws, errWs := db.GetUserWorkspaceByID(currentUser.ID, workspaceID)
		if errWs != nil {
			logrus.Errorln("could not retrieve workspace", errWs)
			return nil, NewGqlError(common.ErrorGeneric, "could not retrieve workspace")
		}
		return *ws, nil
	}

	if err.ErrorCode == common.ErrorDuplicateEntity {
		return nil, NewGqlError(common.ErrorDuplicateEntity, "mock service already added this workspace")
	}

	return nil, NewGqlError(common.ErrorGeneric, "could not add mock service to workspace")
}
