package common

import (
	"os"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
)

func getUserFromClaims(c *jwt.MapClaims) (*models.User, error) {
	u := &models.User{}
	err := mapstructure.Decode(c, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func ParseAuthToken(c *gin.Context) {
	c.Set("user", &models.User{})
	cookie, err := c.Cookie(os.Getenv("AUTH_COOKIE_NAME"))
	if err != nil {
		return
	}
	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH_TOKEN_KEY")), nil
	})

	user, err := getUserFromClaims(&claims)
	if err == nil {
		c.Set("user", user)
	}
}
