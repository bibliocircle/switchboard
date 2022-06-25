package management_api

import (
	"io"
	"log"
	"switchboard/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
	r.GET("/mockservices", GetMockServicesRoute)
	r.GET("/mockservice/:mockServiceId/upstreams", GetUpstreamsByMockServiceIdRoute)
	r.GET("/mockservice/:mockServiceId/endpoints", GetEndpointsByMockServiceIdRoute)

	r.POST("/workspace", CreateWorkspaceRoute)
	r.DELETE("/workspace/:workspaceId", DeleteWorkspaceRoute)
	r.GET("/workspaces", GetWorkspacesRoute)
	r.GET("/user/workspaces", GetUserWorkspacesRoute)
	r.GET("/workspace/:workspaceId/settings", GetWorkspaceSettingsRoute)
	r.PUT("/workspace/:workspaceId/mockservice/:mockServiceId/endpoint/:endpointId/settings", UpdateWsMockServiceConfigRoute)
	r.PUT("/workspace/:workspaceId/mockservice/:mockServiceId/endpoint/:endpointId/scenario/:scenarioId/activate", ActivateMockServiceScenarioRoute)
	r.POST("/workspace/:workspaceId/mockservice/:mockServiceId/add", AddMockServiceToWorkspaceRoute)

	// temporary endpoints to test random data generator
	r.POST("/randomjson", func(c *gin.Context) {
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(500)
		}
		c.Header("Content-Type", "application/json")

		c.Stream(func(w io.Writer) bool {
			err := common.GenFakeJson(string(jsonData), w)
			if err != nil {
				log.Println(err)
			}
			return false
		})
	})
}

func CreateRouter(name string) *gin.Engine {
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
	return r
}
