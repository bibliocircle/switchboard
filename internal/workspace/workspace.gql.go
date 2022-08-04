package workspace

import (
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/gql"
	"switchboard/internal/user"

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
		"expiresAt": &graphql.Field{
			Type: graphql.String,
		},
		"createdBy": &graphql.Field{
			Type: user.UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(Workspace).CreatedBy
				users, err := user.GetUserByID(userId)
				if err != nil {
					logrus.Errorln(err)
					return make([]user.User, 0), gql.NewGqlError(common.ErrorGeneric, "could not resolve createdBy field")
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
	var input CreateWorkspaceRequestBody
	mapstructure.Decode(p.Args["workspace"], &input)

	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	createdWs, createErr := CreateWorkspace(currentUser.ID, &input)

	switch {
	case createErr == nil:
		return createdWs, nil
	case createErr.ErrorCode == common.ErrorInvalidInput:
		return nil, gql.NewGqlError(common.ErrorInvalidInput, "could not parse input")
	case createErr.ErrorCode == common.ErrorDuplicateEntity:
		return nil, gql.NewGqlError(common.ErrorDuplicateEntity, "workspace already exists")
	}

	return nil, gql.NewGqlError(common.ErrorGeneric, "could not create workspace")
}

func GetWorkspaceResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	user := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	if ok {
		ws, err := GetUserWorkspaceByID(user.ID, id)
		if err != nil {
			return Workspace{}, gql.NewGqlError(common.ErrorGeneric, "could not retrieve workspace")
		}
		return *ws, nil
	}
	return Workspace{}, nil
}

func GetWorkspacesResolver(p graphql.ResolveParams) (interface{}, error) {
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	wss, err := GetWorkspaces(currentUser.ID)
	if err != nil {
		logrus.Errorln(err)
		return make([]Workspace, 0), gql.NewGqlError(common.ErrorGeneric, "could not retrieve user workspaces")
	}
	return wss, nil
}

func DeleteWorkspaceResolver(p graphql.ResolveParams) (interface{}, error) {
	workspaceID := p.Args["workspaceId"].(string)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	ok, err := DeleteWorkspace(currentUser.ID, workspaceID)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("could not delete workspace %s. [error code: %s] [description: %s]", workspaceID, err.ErrorCode, err.Description))
		return false, gql.NewGqlError(common.ErrorGeneric, "could not delete workspace")
	}
	if !ok {
		return false, gql.NewGqlError(common.ErrorNotFound, "workspace not found")
	}

	return true, nil
}
