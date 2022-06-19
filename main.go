package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"switchboard/internal/common/middleware"
	"switchboard/internal/common/randomdata"
	"switchboard/internal/common/validation"
	"switchboard/internal/db"
	mapi "switchboard/internal/management_api"

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

	r.POST("/auth/login", mapi.LoginRoute)
	r.POST("/auth/signup", mapi.SignUpRoute)
	r.POST("/auth/logout", mapi.LogOutRoute)
}

func setupAuthenticatedRoutes(r *gin.Engine) {
	r.POST("/endpoint", mapi.CreateEndpointRoute)
	r.DELETE("/endpoint/:endpointId", mapi.DeleteEndpointRoute)

	r.POST("/scenario", mapi.CreateScenarioRoute)

	r.POST("/upstream", mapi.CreateUpstreamRoute)
	r.DELETE("/upstream/:upstreamId", mapi.DeleteUpstreamRoute)

	r.POST("/mockservice", mapi.CreateMockServiceRoute)
	r.DELETE("/mockservice/:mockServiceId", mapi.DeleteMockServiceRoute)
	r.GET("/mockservices", mapi.GetMockServicesRoute)
	r.GET("/mockservice/:mockServiceId/upstreams", mapi.GetUpstreamsByMockServiceIdRoute)
	r.GET("/mockservice/:mockServiceId/endpoints", mapi.GetEndpointsByMockServiceIdRoute)

	r.POST("/workspace", mapi.CreateWorkspaceRoute)
	r.DELETE("/workspace/:workspaceId", mapi.DeleteWorkspaceRoute)
	r.GET("/workspaces", mapi.GetWorkspacesRoute)
	r.GET("/user/workspaces", mapi.GetUserWorkspacesRoute)
	r.GET("/workspace/:workspaceId/settings", mapi.GetWorkspaceSettingsRoute)
	r.PUT("/workspace/:workspaceId/mockservice/:mockServiceId/endpoint/:endpointId/settings", mapi.UpdateWsMockServiceConfigRoute)
	r.PUT("/workspace/:workspaceId/mockservice/:mockServiceId/endpoint/:endpointId/scenario/:scenarioId/activate", mapi.ActivateMockServiceScenarioRoute)
	r.POST("/workspace/:workspaceId/mockservice/:mockServiceId/add", mapi.AddMockServiceToWorkspaceRoute)

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

	r.Use(middleware.ParseAuthToken())
	r.Use(middleware.RequireAuthentication())

	setupAuthenticatedRoutes(r)

	log.Println("Starting server...")
	log.Fatal(r.Run(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
