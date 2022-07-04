package consumer_api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createMiddleware(rc *RouteConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// configure default response headers
		for k, v := range rc.MockService.Config.InjectHeaders {
			c.Header(k, v)
		}
		c.Next()
	}
}

func InitRoute(r *gin.Engine, cfg RouteConfig) {
	routePath := fmt.Sprintf("/ws/:workspaceId/%s%s", cfg.MockService.ID, cfg.Endpoint.Path)
	middleware := createMiddleware(&cfg)
	routeHandler := func(c *gin.Context) {
		wsID := c.Param("workspaceId")
		msID := cfg.MockService.ID
		path := cfg.Endpoint.Path

		c.JSON(http.StatusAccepted, gin.H{
			"wsID": wsID,
			"msID": msID,
			"path": path,
		})
	}

	r.Handle(cfg.Endpoint.Method, routePath, middleware, routeHandler)
	fmt.Printf("initialised endpoint %s %s\n", cfg.Endpoint.Method, routePath)
}
