package gql

import (
	"github.com/graphql-go/graphql"
)

var RootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"users": &graphql.Field{
				Type:    graphql.NewList(UserGqlType),
				Resolve: GetUsersResolver,
			},
			"user": &graphql.Field{
				Type: UserGqlType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: GetUserResolver,
			},
			"mockServices": &graphql.Field{
				Type:    graphql.NewList(MockServiceGqlType),
				Resolve: GetMockServicesResolver,
			},
			"mockService": &graphql.Field{
				Type: MockServiceGqlType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: GetMockServiceResolver,
			},
			"workspaces": &graphql.Field{
				Type:    graphql.NewList(WorkspaceGqlType),
				Resolve: GetWorkspacesResolver,
			},
			"userWorkspaces": &graphql.Field{
				Type:    graphql.NewList(WorkspaceGqlType),
				Resolve: GetUserWorkspacesResolver,
			},
			"userWorkspace": &graphql.Field{
				Type: WorkspaceGqlType,
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: GetUserWorkspaceResolver,
			},
			"workspaceSettings": &graphql.Field{
				Type: graphql.NewList(WorkspaceSettingGqlType),
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: GetWorkspaceSettingsResolver,
			},
			"workspaceSetting": &graphql.Field{
				Type: WorkspaceSettingGqlType,
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"mockServiceId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: GetWorkspaceSettingResolver,
			},
		},
	},
)

var RootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createWorkspace": &graphql.Field{
				Type: WorkspaceGqlType,
				Args: graphql.FieldConfigArgument{
					"workspace": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(WorkspaceGqlInputType),
					},
				},
				Resolve: CreateWorkspaceResolver,
			},
			"deleteWorkspace": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: DeleteWorkspaceResolver,
			},
			"createEndpoint": &graphql.Field{
				Type: EndpointGqlType,
				Args: graphql.FieldConfigArgument{
					"endpoint": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(EndpointGqlInputType),
					},
				},
				Resolve: CreateEndpointResolver,
			},
			"deleteEndpoint": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"endpointId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: DeleteEndpointResolver,
			},
			"createScenario": &graphql.Field{
				Type: ScenarioGqlType,
				Args: graphql.FieldConfigArgument{
					"scenario": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(ScenarioGqlInputType),
					},
				},
				Resolve: CreateScenarioResolver,
			},
			"createUpstream": &graphql.Field{
				Type: UpstreamGqlType,
				Args: graphql.FieldConfigArgument{
					"upstream": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(UpstreamGqlInputType),
					},
				},
				Resolve: CreateUpstreamResolver,
			},
			"deleteUpstream": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"upstreamId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: DeleteUpstreamResolver,
			},
			"createMockService": &graphql.Field{
				Type: MockServiceGqlType,
				Args: graphql.FieldConfigArgument{
					"mockService": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(MockServiceGqlInputType),
					},
				},
				Resolve: CreateMockServiceResolver,
			},
			"deleteMockService": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"mockServiceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: DeleteMockServiceResolver,
			},
			"activateMockServiceScenario": &graphql.Field{
				Type: WorkspaceSettingGqlType,
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
				Resolve: ActivateMockServiceScenarioResolver,
			},
			"addMockServiceToWorkspace": &graphql.Field{
				Type: WorkspaceGqlType,
				Args: graphql.FieldConfigArgument{
					"workspaceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"mockServiceId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: AddMockServiceToWorkspaceResolver,
			},
		},
	},
)
