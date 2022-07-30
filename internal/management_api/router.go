package management_api

import (
	"io"
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/gql"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/graph-gophers/dataloader"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func setupUnauthenticatedRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.Writer.Write([]byte("pong"))
	})

	r.POST("/auth/login", LoginRoute)
	r.POST("/auth/signup", SignUpRoute)
	r.POST("/auth/logout", LogOutRoute)
}

func setupAuthenticatedRoutes(r *gin.Engine) {
	r.POST("/endpoint", CreateEndpointRoute)
	r.DELETE("/endpoint/:endpointId", DeleteEndpointRoute)

	r.POST("/scenario", CreateScenarioRoute)

	r.POST("/upstream", CreateUpstreamRoute)
	r.DELETE("/upstream/:upstreamId", DeleteUpstreamRoute)

	r.POST("/mockservice", CreateMockServiceRoute)
	r.DELETE("/mockservice/:mockServiceId", DeleteMockServiceRoute)

	r.POST("/workspace", CreateWorkspaceRoute)
	r.DELETE("/workspace/:workspaceId", DeleteWorkspaceRoute)

	// temporary endpoints to test random data generator
	r.POST("/randomjson", func(c *gin.Context) {
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		data, fErr := common.GenFakeData(string(jsonData))
		if fErr != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		c.Header("Content-Type", "application/json")
		c.Writer.Write([]byte(data))
	})
}

func CreateRouter(name string, reload chan bool, quit chan<- bool) *gin.Engine {
	schemaConfig := graphql.SchemaConfig{
		Query:    gql.RootQuery,
		Mutation: gql.RootMutation,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		common.InitialiseValidator(v)
	}

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: common.CreateGinLogFormatter(name),
	}))
	r.Use(gin.Recovery())
	r.Use(common.ConfigureCors())

	setupUnauthenticatedRoutes(r)

	r.Use(common.ParseAuthToken())
	r.Use(common.RequireAuthentication())

	setupAuthenticatedRoutes(r)
	r.Any("/graphql", func(ctx *gin.Context) {
		loaders := &db.Loaders{
			Scenarios: dataloader.NewBatchedLoader(db.BatchLoadScenarios),
			Endpoints: dataloader.NewBatchedLoader(db.BatchLoadEndpoints),
		}
		ctx.Set(db.LoadersCtxKey, loaders)
		h.ServeHTTP(ctx.Writer, ctx.Request.WithContext(ctx))
	})
	return r
}
