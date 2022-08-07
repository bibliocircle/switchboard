package workspace_setting

import (
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/endpoint"
	"switchboard/internal/gql"
	"switchboard/internal/mockservice"
	"switchboard/internal/scenario"
	"switchboard/internal/user"
	"switchboard/internal/workspace"

	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var InterceptionRuleGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "InterceptionRuleInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"matcherExpression": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"TargetScenarioId": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var InterceptionRuleGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "InterceptionRule",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"matcherExpression": &graphql.Field{
			Type: graphql.String,
		},
		"targetScenarioId": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ScenarioConfig",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"scenario": &graphql.Field{
			Type: scenario.ScenarioGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				scenarioID := p.Source.(ScenarioConfig).ScenarioID
				loaders := p.Context.Value(db.LoadersCtxKey).(*db.Loaders)
				thunk := loaders.Scenarios.Load(p.Context, dataloader.StringKey(scenarioID))
				return func() (interface{}, error) {
					sc, err := thunk()
					if err != nil {
						logrus.Errorln(err)
						return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve scenario")
					}
					return sc.(*scenario.Scenario), nil
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
			Type:    endpoint.EndpointGqlType,
			Resolve: GetEndpointResolver,
		},
		"scenarioConfigs": &graphql.Field{
			Type: graphql.NewList(ScenarioConfigGqlType),
		},
		"interceptionRules": &graphql.Field{
			Type: graphql.NewList(InterceptionRuleGqlType),
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
			Type: workspace.WorkspaceGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				workspaceID := p.Source.(*WorkspaceSetting).WorkspaceID
				currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
				ws, err := workspace.GetUserWorkspaceByID(currentUser.ID, workspaceID)
				if err != nil {
					logrus.Errorln(err)
					return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve workspace")
				}
				return ws, nil
			},
		},
		"mockService": &graphql.Field{
			Type: mockservice.MockServiceGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mockServiceID := p.Source.(*WorkspaceSetting).MockServiceID
				ms, err := mockservice.GetMockServiceByID(mockServiceID)
				if err != nil {
					logrus.Errorln(err)
					return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve mock service")
				}
				return ms, nil
			},
		},
		"config": &graphql.Field{
			Type: mockservice.GlobalMockServiceConfigGqlType,
		},
		"endpointConfigs": &graphql.Field{
			Type: graphql.NewList(EndpointConfigGqlType),
		},
	},
})

func GetEndpointResolver(p graphql.ResolveParams) (interface{}, error) {
	endpointID := p.Source.(EndpointConfig).EndpointID
	loaders := p.Context.Value(db.LoadersCtxKey).(*db.Loaders)
	thunk := loaders.Endpoints.Load(p.Context, dataloader.StringKey(endpointID))
	return func() (interface{}, error) {
		ep, err := thunk()
		if err != nil {
			logrus.Errorln(err)
			return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve endpoint")
		}
		return ep.(*endpoint.Endpoint), nil
	}, nil
}

func GetWorkspaceSettingsResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	user := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	if isWsOwner, err := workspace.IsWorkspaceOwner(user.ID, workspaceID); err != nil {
		return []WorkspaceSetting{}, err
	} else if !isWsOwner {
		return []WorkspaceSetting{}, gql.NewGqlError(common.ErrorForbidden, "could not find workspace or the user does not have permission to this workspace")
	}

	wss, err := GetWorkspaceSettings(workspaceID)
	if err != nil {
		logrus.Errorf("could not retrieve workspace settings for workspace ID %s : %s\n", workspaceID, err)
		return []WorkspaceSetting{}, gql.NewGqlError(common.ErrorGeneric, "could not retrieve workspace settings")
	}

	return *wss, nil
}

func GetWorkspaceSettingResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	mockServiceID := p.Args["mockServiceId"].(string)
	wss, err := GetWorkspaceSetting(workspaceID, mockServiceID)
	if err == nil {
		return wss, nil
	}
	if err.ErrorCode == common.ErrorNotFound {
		return nil, gql.NewGqlError(common.ErrorNotFound, "workspace not found or mock service not found in the workspace")
	}
	return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve workspace settings")
}

func ActivateMockServiceScenarioResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	mockServiceID := p.Args["mockServiceId"].(string)
	endpointID := p.Args["endpointId"].(string)
	scenarioID := p.Args["scenarioId"].(string)

	ok, errAct := ActivateWsMockServiceScenario(
		workspaceID,
		mockServiceID,
		endpointID,
		scenarioID,
	)
	if errAct != nil {
		return nil, gql.NewGqlError(common.ErrorGeneric, "could not activate mock service scenario")
	}
	if ok {
		wss, err := GetWorkspaceSetting(workspaceID, mockServiceID)
		if err != nil {
			return nil, err
		}
		var endpointConfig EndpointConfig
		for _, ec := range wss.EndpointConfigs {
			if ec.EndpointID == endpointID {
				endpointConfig = ec
				break
			}
		}
		return endpointConfig, nil
	}
	return nil, nil
}

func AddMockServiceToWorkspaceResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	mockServiceID := p.Args["mockServiceId"].(string)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)

	isWsOwner, errPerm := workspace.IsWorkspaceOwner(currentUser.ID, workspaceID)
	if errPerm != nil {
		logrus.Errorln("could not check workspace permissions", errPerm)
		return nil, gql.NewGqlError(common.ErrorGeneric, "could not check workspace permissions")
	}

	if !isWsOwner {
		return nil, gql.NewGqlError(common.ErrorForbidden, "workspace doesn't exist or the current user doesn't have permissions to access the workspace")
	}

	err := AddMockServiceToWorkspace(currentUser.ID, workspaceID, mockServiceID)
	if err == nil {
		ws, errWs := workspace.GetUserWorkspaceByID(currentUser.ID, workspaceID)
		if errWs != nil {
			logrus.Errorln("could not retrieve workspace", errWs)
			return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve workspace")
		}
		return *ws, nil
	}

	if err.ErrorCode == common.ErrorDuplicateEntity {
		return nil, gql.NewGqlError(common.ErrorDuplicateEntity, "mock service already added this workspace")
	}

	return nil, gql.NewGqlError(common.ErrorGeneric, "could not add mock service to workspace")
}

func CreateInterceptionRuleResolver(p graphql.ResolveParams) (interface{}, error) {
	var interceptionRule CreateInterceptionRuleRequestBody
	workspaceID := p.Args["workspaceId"].(string)
	mockServiceID := p.Args["mockServiceId"].(string)
	endpointID := p.Args["endpointId"].(string)
	mapstructure.Decode(p.Args["interceptionRule"], &interceptionRule)

	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	created, err := CreateInterceptionRule(currentUser.ID, workspaceID, mockServiceID, endpointID, interceptionRule)
	if err != nil {
		logrus.Error(fmt.Sprintf("could not create interception rule for %s\n", endpointID))
		return nil, gql.NewGqlError(common.ErrorGeneric, "could not create interception rule")
	}
	if !created {
		return nil, gql.NewGqlError(common.ErrorNotFound, "no such endpoint found")
	}

	return created, nil
}
