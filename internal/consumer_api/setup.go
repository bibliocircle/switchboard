package consumer_api

import "github.com/gin-gonic/gin"

func CreateRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	return r
}
