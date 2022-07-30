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
