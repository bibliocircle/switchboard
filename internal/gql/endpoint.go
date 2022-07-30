package gql

import (
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var EndpointGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Endpoint",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"path": &graphql.Field{
			Type: graphql.String,
		},
		"method": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"responseDelay": &graphql.Field{
			Type: graphql.Int,
		},
		"scenarios": &graphql.Field{
			Type: graphql.NewList(ScenarioGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				endpointID := p.Source.(models.Endpoint).ID
				scenarios, err := db.GetScenarios(endpointID)
				if err != nil {
					logrus.Errorln(err)
					return make([]models.Scenario, 0), NewGqlError(GqlInternalError, "could not retrieve scenarios")
				}
				return scenarios, nil
			},
		},
		"createdBy": &graphql.Field{
			Type: UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(models.Endpoint).CreatedBy
				users, err := db.GetUserByID(userId)
				if err != nil {
					logrus.Errorln(err)
					return make([]models.User, 0), NewGqlError(GqlInternalError, "could not resolve createdBy field")
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
