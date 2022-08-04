package mockservice

import (
	"fmt"
	"switchboard/internal/common"
	"switchboard/internal/endpoint"
	"switchboard/internal/gql"
	"switchboard/internal/upstream"
	"switchboard/internal/user"

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
			Type: graphql.NewList(gql.HTTPHeaderGqlInputType),
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
			Type: graphql.NewList(gql.HTTPHeaderGqlType),
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
			Type: graphql.NewList(upstream.UpstreamGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mockServiceID := p.Source.(*MockService).ID
				upstreams, err := upstream.GetUpstreams(mockServiceID)
				if err != nil {
					logrus.Errorln(err)
					return make([]upstream.Upstream, 0), gql.NewGqlError(common.ErrorGeneric, "could not retrieve upstreams")
				}
				return upstreams, nil
			},
		},
		"endpoints": &graphql.Field{
			Type: graphql.NewList(endpoint.EndpointGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mockServiceID := p.Source.(*MockService).ID
				endpoints, err := endpoint.GetEndpoints(mockServiceID)
				if err != nil {
					logrus.Errorln(err)
					return make([]endpoint.Endpoint, 0), gql.NewGqlError(common.ErrorGeneric, "could not retrieve endpoints")
				}
				return endpoints, nil
			},
		},
		"createdBy": &graphql.Field{
			Type: user.UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(*MockService).CreatedBy
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

func GetMockServicesResolver(p graphql.ResolveParams) (interface{}, error) {
	svcs, err := GetMockServices()
	if err != nil {
		logrus.Errorln(err)
		return make([]MockService, 0), gql.NewGqlError(common.ErrorGeneric, "could not retrieve mock services")
	}
	return svcs, nil
}

func GetMockServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if ok {
		mockService, err := GetMockServiceByID(id)
		if err != nil {
			return nil, gql.NewGqlError(common.ErrorGeneric, "could not retrieve mock service")
		}
		return mockService, nil
	}
	return nil, nil
}

func CreateMockServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	var input CreateMockServiceRequestBody
	mapstructure.Decode(p.Args["mockService"], &input)

	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	createdMockService, createErr := CreateMockService(currentUser.ID, &input)
	if createErr == nil {
		return createdMockService, nil
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		return nil, gql.NewGqlError(common.ErrorDuplicateEntity, "duplicate mock service")
	}

	return nil, gql.NewGqlError(common.ErrorGeneric, "could not create mock service")
}

func DeleteMockServiceResolver(p graphql.ResolveParams) (interface{}, error) {
	mockServiceID := p.Args["mockServiceId"].(string)
	currentUser := p.Context.Value(common.REQ_USER_KEY).(*user.User)
	ok, err := DeleteMockService(currentUser.ID, mockServiceID)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("could not delete mock service %s. [error code: %s] [description: %s]", mockServiceID, err.ErrorCode, err.Description))
		return false, gql.NewGqlError(common.ErrorGeneric, "could not delete mock service")
	}
	if !ok {
		return false, gql.NewGqlError(common.ErrorNotFound, "mock service not found")
	}
	return true, nil
}
