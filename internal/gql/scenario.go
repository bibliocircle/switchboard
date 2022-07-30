package gql

import (
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var HTTPResponseScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "HTTPResponseScenarioConfig",
	Fields: graphql.Fields{
		"statusCode": &graphql.Field{
			Type: graphql.Int,
		},
		"responseBodyTemplate": &graphql.Field{
			Type: graphql.String,
		},
		"responseHeadersTemplate": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ProxyScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProxyScenarioConfig",
	Fields: graphql.Fields{
		"upstream": &graphql.Field{
			Type: UpstreamGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				upstreamID := p.Source.(*models.ProxyScenarioConfig).UpstreamID
				if upstreamID == "" {
					return nil, nil
				}
				upstream, err := db.GetUpstreamByID(upstreamID)
				if err != nil {
					logrus.Errorln(err)
					return nil, NewGqlError(GqlInternalError, "could not retrieve upstream")
				}
				return *upstream, nil
			},
		},
		"injectHeaders": &graphql.Field{
			Type: graphql.NewList(HTTPHeaderGqlType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				headers := make([]HTTPHeader, 0)
				injectHeaders := p.Source.(*models.ProxyScenarioConfig).InjectHeaders
				for k, v := range injectHeaders {
					headers = append(headers, HTTPHeader{
						Name:  k,
						Value: v,
					})
				}
				return headers, nil
			},
		},
	},
})

var NetworkScenarioConfigGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "NetworkScenarioConfig",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var ScenarioGqlType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Scenario",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"httpResponseScenarioConfig": &graphql.Field{
			Type: HTTPResponseScenarioConfigGqlType,
		},
		"proxyScenarioConfig": &graphql.Field{
			Type: ProxyScenarioConfigGqlType,
		},
		"networkScenarioConfig": &graphql.Field{
			Type: NetworkScenarioConfigGqlType,
		},
		"createdBy": &graphql.Field{
			Type: UserGqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				userId := p.Source.(models.Scenario).CreatedBy
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
