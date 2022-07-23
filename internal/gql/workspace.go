package gql

import (
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
)

var WorkspaceGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Workspace",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"expiresAt": &graphql.Field{
			Type: graphql.String,
		},
		"createdBy": &graphql.Field{
			Type: UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(models.Workspace).CreatedBy
				users, err := db.GetUserByID(userId)
				if err != nil {
					return make([]models.User, 0), err
				}
				return users, nil
			},
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var WorkspacesResolver = func(p graphql.ResolveParams) (interface{}, error) {
	wss, err := db.GetWorkspaces()
	if err != nil {
		return make([]models.Workspace, 0), err
	}
	return wss, nil
}
