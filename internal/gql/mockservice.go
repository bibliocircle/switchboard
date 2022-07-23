package gql

import (
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
)

var GlobalMockServiceConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "GlobalMockServiceConfig",
	Fields: graphql.Fields{
		"injectHeaders": &graphql.Field{
			Type: graphql.NewList(HeaderConfigGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				headers := make([]HeaderConfig, 0)
				injectHeaders := p.Source.(models.GlobalMockServiceConfig).InjectHeaders
				for k, v := range injectHeaders {
					headers = append(headers, HeaderConfig{
						Name:  k,
						Value: v,
					})
				}
				return headers, nil
			},
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
				mockServiceID := p.Source.(models.MockService).ID
				upstreams, err := db.GetUpstreams(mockServiceID)
				if err != nil {
					return make([]models.Upstream, 0), err
				}
				return upstreams, nil
			},
		},
		"endpoints": &graphql.Field{
			Type: graphql.NewList(EndpointGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				mockServiceID := p.Source.(models.MockService).ID
				endpoints, err := db.GetEndpoints(mockServiceID)
				if err != nil {
					return make([]models.Endpoint, 0), err
				}
				return endpoints, nil
			},
		},
		"createdBy": &graphql.Field{
			Type: UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(models.MockService).CreatedBy
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

var MockServicesResolver = func(p graphql.ResolveParams) (interface{}, error) {
	svcs, err := db.GetMockServices()
	if err != nil {
		return make([]models.MockService, 0), err
	}
	return svcs, nil
}

var MockServiceResolver = func(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if ok {
		mockService, err := db.GetMockServiceByID(id)
		if err != nil {
			return nil, err
		}
		return *mockService, nil
	}
	return nil, nil
}
