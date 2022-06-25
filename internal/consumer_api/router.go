package consumer_api

import (
	"net/http"
	"switchboard/internal/common"

	"github.com/gin-gonic/gin"
)

func CreateRouter(name string) *gin.Engine {
	r := gin.New()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: common.CreateGinLogFormatter(name),
	}))
	r.Use(gin.Recovery())

	/*
	* Request format: GET http://localhost:9999/ws/1234/order-api/order/20
	 */
	r.Any("/ws/:workspaceId/:mockServiceId/*path", func(c *gin.Context) {
		wsID := c.Param("workspaceId")
		msID := c.Param("mockServiceId")
		path := c.Param("path")

		c.JSON(http.StatusAccepted, gin.H{
			"wsID": wsID,
			"msID": msID,
			"path": path,
		})
	})

	return r
}
