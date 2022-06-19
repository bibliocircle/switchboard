package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/db"
	"switchboard/internal/common/middleware"
	"switchboard/internal/common/randomdata"
	"switchboard/internal/common/validation"
	"switchboard/internal/management_api/endpoint"
	"switchboard/internal/management_api/mockservice"
	"switchboard/internal/management_api/scenario"
	"switchboard/internal/management_api/upstream"
	"switchboard/internal/management_api/workspace"
	"switchboard/internal/management_api/workspace_settings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	env "github.com/joho/godotenv"
)

func setUpDatabase() {
	ctx := context.Background()
	dbError := db.Connect(ctx)
	if dbError != nil {
		log.Fatalln("could not connect to the database", dbError)
	}
	db.Migrate(db.Client)
}

func setupUnauthenticatedRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.Writer.Write([]byte("pong"))
	})

	r.POST("/auth/login", auth.LoginRoute)
	r.POST("/auth/signup", auth.SignUpRoute)
	r.POST("/auth/logout", auth.LogOutRoute)
}

func setupAuthenticatedRoutes(r *gin.Engine) {
	r.POST("/endpoint", endpoint.CreateEndpointRoute)
	r.DELETE("/endpoint/:endpointId", endpoint.DeleteEndpointRoute)

	r.POST("/scenario", scenario.CreateScenarioRoute)

	r.POST("/upstream", upstream.CreateUpstreamRoute)
	r.DELETE("/upstream/:upstreamId", upstream.DeleteUpstreamRoute)

	r.POST("/mockservice", mockservice.CreateMockServiceRoute)
	r.DELETE("/mockservice/:mockServiceId", mockservice.DeleteMockServiceRoute)
	r.GET("/mockservices", mockservice.GetMockServicesRoute)
	r.GET("/mockservice/:mockServiceId/upstreams", upstream.GetUpstreamsByMockServiceIdRoute)
	r.GET("/mockservice/:mockServiceId/endpoints", endpoint.GetEndpointsByMockServiceIdRoute)

	r.POST("/workspace", workspace.CreateWorkspaceRoute)
	r.DELETE("/workspace/:workspaceId", workspace.DeleteWorkspaceRoute)
	r.GET("/workspaces", workspace.GetWorkspacesRoute)
	r.GET("/user/workspaces", workspace.GetUserWorkspacesRoute)
	r.GET("/workspace/:workspaceId/settings", workspace_settings.GetWorkspaceSettingsRoute)
	r.PUT("/workspace/:workspaceId/mockservice/:mockServiceId/endpoint/:endpointId/settings", workspace_settings.UpdateWsMockServiceConfigRoute)
	r.PUT("/workspace/:workspaceId/mockservice/:mockServiceId/endpoint/:endpointId/scenario/:scenarioId/activate", workspace_settings.ActivateMockServiceScenarioRoute)
	r.POST("/workspace/:workspaceId/mockservice/:mockServiceId/add", workspace_settings.AddMockServiceToWorkspaceRoute)

	// temporary endpoints to test random data generator
	r.POST("/randomjson", func(c *gin.Context) {
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(500)
		}
		c.Header("Content-Type", "application/json")

		c.Stream(func(w io.Writer) bool {
			err := randomdata.GenFakeJson(string(jsonData), w)
			if err != nil {
				log.Println(err)
			}
			return false
		})
	})
}

func main() {
	err := env.Load()
	if err != nil {
		log.Println("could not locate or read .env file", err)
	}

	setUpDatabase()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.InitialiseValidator(v)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ConfigureCors())

	setupUnauthenticatedRoutes(r)

	r.Use(auth.ParseAuthToken())
	r.Use(auth.RequireAuthentication())

	setupAuthenticatedRoutes(r)

	log.Println("Starting server...")
	log.Fatal(r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
