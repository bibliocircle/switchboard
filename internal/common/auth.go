package common

import (
	"net/http"
	"os"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func getUserFromClaims(c *jwt.MapClaims) (*models.User, error) {
	u := &models.User{}
	err := mapstructure.Decode(c, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func ParseAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(os.Getenv("AUTH_COOKIE_NAME"))
		if err != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Abort()
			return
		}
		claims := jwt.MapClaims{}
		_, tErr := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("AUTH_TOKEN_KEY")), nil
		})

		if tErr != nil {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Abort()
			return
		}
		user, err := getUserFromClaims(&claims)
		if err == nil {
			c.Set("user", user)
		} else {
			log.Error("could not extract claims from auth token")
		}
	}
}

func RequireAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := c.Value("user").(*models.User)
		if currentUser.ID == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}
