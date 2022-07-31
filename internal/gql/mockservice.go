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

var MockServiceGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "MockServiceInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"type": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"config": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(GlobalMockServiceConfigGqlInputType),
		},
	},
})

var GlobalMockServiceConfigGqlInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "GlobalMockServiceConfigInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"injectHeaders": &graphql.InputObjectFieldConfig{
			Type: graphql.NewList(HTTPHeaderGqlInputType),
		},
		"globalResponseDelay": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
	},
})

var GlobalMockServiceConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GlobalMockServiceConfig",
	Fields: graphql.Fields{
		"injectHeaders": &graphql.Field{
			Type: graphql.NewList(HTTPHeaderGqlType),
		},
		"globalResponseDelay": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var MockServiceGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "MockService",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"config": &graphql.Field{
			Type: GlobalMockServiceConfigGqlType,
		},
		"upstreams": &graphql.Field{
			Type: graphql.NewList(UpstreamGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mockServiceID := p.Source.(*models.MockService).ID
				upstreams, err := db.GetUpstreams(mockServiceID)
				if err != nil {
					logrus.Errorln(err)
					return make([]models.Upstream, 0), NewGqlError(common.ErrorGeneric, "could not retrieve upstreams")
				}
				return upstreams, nil
			},
		},
		"endpoints": &graphql.Field{
			Type: graphql.NewList(EndpointGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mockServiceID := p.Source.(*models.MockService).ID
				endpoints, err := db.GetEndpoints(mockServiceID)
				if err != nil {
					logrus.Errorln(err)
					return make([]models.Endpoint, 0), NewGqlError(common.ErrorGeneric, "could not retrieve endpoints")
				}
				return endpoints, nil
			},
		},
		"createdBy": &graphql.Field{
			Type: UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(*models.MockService).CreatedBy
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

func GetMockServicesResolver(p graphql.ResolveParams) (interface{}, error) {
	svcs, err := db.GetMockServices()
	if err != nil {
		logrus.Errorln(err)
		return make([]models.MockService, 0), NewGqlError(common.ErrorGeneric, "could not retrieve mock services")
	}
	return svcs, nil
}

func GetMockServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if ok {
		mockService, err := db.GetMockServiceByID(id)
		if err != nil {
			return nil, NewGqlError(common.ErrorGeneric, "could not retrieve mock service")
		}
		return mockService, nil
	}
	return nil, nil
}

func CreateMockServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	var input models.CreateMockServiceRequestBody
	mapstructure.Decode(p.Args["mockService"], &input)

	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
	createdMockService, createErr := db.CreateMockService(currentUser.ID, &input)
	if createErr == nil {
		return createdMockService, nil
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		return nil, NewGqlError(common.ErrorDuplicateEntity, "duplicate mock service")
	}

	return nil, NewGqlError(common.ErrorGeneric, "could not create mock service")
}

func DeleteMockServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	mockServiceID := p.Args["mockServiceId"].(string)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*models.User)
	ok, err := db.DeleteMockService(currentUser.ID, mockServiceID)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("could not delete mock service %s. [error code: %s] [description: %s]", mockServiceID, err.ErrorCode, err.Description))
		return false, NewGqlError(common.ErrorGeneric, "could not delete mock service")
	}
	if !ok {
		return false, NewGqlError(common.ErrorNotFound, "mock service not found")
	}
	return true, nil
}
