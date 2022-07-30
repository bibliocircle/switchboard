package gql

import (
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

var WorkspaceGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "WorkspaceInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"expiresAt": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

var WorkspaceGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Workspace",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"mockServices": &graphql.Field{
			Type:    graphql.NewList(MockServiceGqlType),
			Resolve: GetWorkspaceMockServicesResolver,
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
					logrus.Errorln(err)
					return make([]models.User, 0), NewGqlError(common.ErrorGeneric, "could not resolve createdBy field")
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

func CreateWorkspaceResolver(p graphql.ResolveParams) (interface{}, error) {
	var input models.CreateWorkspaceRequestBody
	mapstructure.Decode(p.Args["workspace"], &input)

	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
	createdWs, createErr := db.CreateWorkspace(currentUser.ID, &input)

	switch {
	case createErr == nil:
		return createdWs, nil
	case createErr.ErrorCode == common.ErrorInvalidInput:
		return nil, NewGqlError(common.ErrorInvalidInput, "could not parse input")
	case createErr.ErrorCode == common.ErrorDuplicateEntity:
		return nil, NewGqlError(common.ErrorDuplicateEntity, "workspace already exists")
	}

	return nil, NewGqlError(common.ErrorGeneric, "could not create workspace")
}

func GetWorkspacesResolver(p graphql.ResolveParams) (interface{}, error) {
	wss, err := db.GetWorkspaces()
	if err != nil {
		logrus.Errorln(err)
		return make([]models.Workspace, 0), NewGqlError(common.ErrorGeneric, "could not retrieve workspaces")
	}
	return wss, nil
}

func GetUserWorkspacesResolver(p graphql.ResolveParams) (interface{}, error) {
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
	wss, err := db.GetUserWorkspaces(currentUser.ID)
	if err != nil {
		logrus.Errorln(err)
		return make([]models.Workspace, 0), NewGqlError(common.ErrorGeneric, "could not retrieve user workspaces")
	}
	return wss, nil
}

func GetUserWorkspaceResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID, ok := p.Args["workspaceId"].(string)
	if ok {
		currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
		wss, err := db.GetUserWorkspaceByID(currentUser.ID, workspaceID)
		if err != nil {
			return nil, err
		}
		return *wss, nil
	}
	return nil, nil
}

func DeleteWorkspaceResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
	ok, err := db.DeleteWorkspace(currentUser.ID, workspaceID)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("could not delete workspace %s. [error code: %s] [description: %s]", workspaceID, err.ErrorCode, err.Description))
		return false, NewGqlError(common.ErrorGeneric, "could not delete workspace")
	}
	if !ok {
		return false, NewGqlError(common.ErrorNotFound, "workspace not found")
	}

	return true, nil
}
