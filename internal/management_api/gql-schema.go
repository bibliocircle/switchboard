package management_api

import (
	"switchboard/internal/endpoint"
	"switchboard/internal/mockservice"
	"switchboard/internal/scenario"
	"switchboard/internal/upstream"
	"switchboard/internal/user"
	"switchboard/internal/workspace"
	"switchboard/internal/workspace_setting"

	"github.com/graphql-go/graphql"
)

var RootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type:    graphql.NewList(user.UserGqlType),
				Resolve: user.GetUsersResolver,
			},
			"user": &graphql.Field{
				Type: user.UserGqlType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: user.GetUserResolver,
			},
			"mockServices": &graphql.Field{
				Type:    graphql.NewList(mockservice.MockServiceGqlType),
				Resolve: mockservice.GetMockServicesResolver,
			},
			"mockService": &graphql.Field{
				Type: mockservice.MockServiceGqlType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: mockservice.GetMockServiceResolver,
			},
			"workspace": &graphql.Field{
				Type: workspace.WorkspaceGqlType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: workspace.GetWorkspaceResolver,
			},
			"workspaces": &graphql.Field{
				Type:    graphql.NewList(workspace.WorkspaceGqlType),
				Resolve: workspace.GetWorkspacesResolver,
			},
			"workspaceSettings": &graphql.Field{
				Type: graphql.NewList(workspace_setting.WorkspaceSettingGqlType),
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: workspace_setting.GetWorkspaceSettingsResolver,
			},
			"workspaceSetting": &graphql.Field{
				Type: workspace_setting.WorkspaceSettingGqlType,
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"mockServiceId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: workspace_setting.GetWorkspaceSettingResolver,
			},
		},
	},
)

var RootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createWorkspace": &graphql.Field{
				Type: workspace.WorkspaceGqlType,
				Args: graphql.FieldConfigArgument{
					"workspace": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(workspace.WorkspaceGqlInputType),
					},
				},
				Resolve: workspace.CreateWorkspaceResolver,
			},
			"deleteWorkspace": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: workspace.DeleteWorkspaceResolver,
			},
			"createEndpoint": &graphql.Field{
				Type: endpoint.EndpointGqlType,
				Args: graphql.FieldConfigArgument{
					"endpoint": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(endpoint.EndpointGqlInputType),
					},
				},
				Resolve: endpoint.CreateEndpointResolver,
			},
			"deleteEndpoint": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"endpointId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: endpoint.DeleteEndpointResolver,
			},
			"createScenario": &graphql.Field{
				Type: scenario.ScenarioGqlType,
				Args: graphql.FieldConfigArgument{
					"scenario": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(scenario.ScenarioGqlInputType),
					},
				},
				Resolve: scenario.CreateScenarioResolver,
			},
			"createUpstream": &graphql.Field{
				Type: upstream.UpstreamGqlType,
				Args: graphql.FieldConfigArgument{
					"upstream": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(upstream.UpstreamGqlInputType),
					},
				},
				Resolve: upstream.CreateUpstreamResolver,
			},
			"deleteUpstream": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"upstreamId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: upstream.DeleteUpstreamResolver,
			},
			"createMockService": &graphql.Field{
				Type: mockservice.MockServiceGqlType,
				Args: graphql.FieldConfigArgument{
					"mockService": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(mockservice.MockServiceGqlInputType),
					},
				},
				Resolve: mockservice.CreateMockServiceResolver,
			},
			"deleteMockService": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"mockServiceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: mockservice.DeleteMockServiceResolver,
			},
			"activateMockServiceScenario": &graphql.Field{
				Type: workspace_setting.EndpointConfigGqlType,
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"mockServiceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"endpointId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"scenarioId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: workspace_setting.ActivateMockServiceScenarioResolver,
			},
			"addMockServiceToWorkspace": &graphql.Field{
				Type: workspace.WorkspaceGqlType,
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"mockServiceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: workspace_setting.AddMockServiceToWorkspaceResolver,
			},
		},
	},
)
