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

var UpstreamGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpstreamInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"mockServiceId": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"url": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var UpstreamGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Upstream",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"createdBy": &graphql.Field{
			Type: UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(models.Upstream).CreatedBy
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

func CreateUpstreamResolver(p graphql.ResolveParams) (interface{}, error) {
	var input models.CreateUpstreamRequestBody
	mapstructure.Decode(p.Args["upstream"], &input)

	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
	createdUpstream, createErr := db.CreateUpstream(currentUser.ID, &input)
	if createErr == nil {
		return createdUpstream, nil
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		return nil, NewGqlError(common.ErrorDuplicateEntity, "duplicate upstream")
	}

	return nil, NewGqlError(common.ErrorGeneric, "could not create upstream")
}

func DeleteUpstreamResolver(p graphql.ResolveParams) (interface{}, error) {
	upstreamID := p.Args["upstreamId"].(string)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
	ok, err := db.DeleteUpstream(currentUser.ID, upstreamID)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("could not delete upstream %s. [error code: %s] [description: %s]", upstreamID, err.ErrorCode, err.Description))
		return false, NewGqlError(common.ErrorGeneric, "could not delete upstream")
	}
	if !ok {
		return false, NewGqlError(common.ErrorNotFound, "upstream not found")
	}

	return true, nil
}
