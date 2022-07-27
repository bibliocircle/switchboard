package gql

import (
	"github.com/graphql-go/graphql"
)

var GqlSchema = graphql.Fields{
	"users": &graphql.Field{
		Type:    graphql.NewList(UserGqlType),
		Resolve: UsersResolver,
	},
	"user": &graphql.Field{
		Type: UserGqlType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: UserResolver,
	},
	"mockServices": &graphql.Field{
		Type:    graphql.NewList(MockServiceGqlType),
		Resolve: MockServicesResolver,
	},
	"mockService": &graphql.Field{
		Type: MockServiceGqlType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: MockServiceResolver,
	},
	"workspaces": &graphql.Field{
		Type:    graphql.NewList(WorkspaceGqlType),
		Resolve: WorkspacesResolver,
	},
	"userWorkspaces": &graphql.Field{
		Type:    graphql.NewList(WorkspaceGqlType),
		Resolve: UserWorkspacesResolver,
	},
	"userWorkspace": &graphql.Field{
		Type: WorkspaceGqlType,
		Args: graphql.FieldConfigArgument{
			"workspaceId": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: UserWorkspaceResolver,
	},
}
