package common

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func isOriginAllowed(origin string) bool {
	allowedOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	for _, url := range allowedOrigins {
		if url == origin {
			return true
		}
	}
	return false
}

func ConfigureCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if isOriginAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
}
