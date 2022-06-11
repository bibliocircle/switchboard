package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type User struct {
	FirstName string `mapstructure:"firstName"`
	Lastname  string `mapstructure:"lastName"`
	Email     string `mapstructure:"email"`
	ID        string `mapstructure:"id"`
	Role      string `mapstructure:"role"`
}

func getUserFromClaims(c *jwt.MapClaims) (*User, error) {
	u := &User{}
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
		currentUser := c.Value("user").(*User)
		if currentUser.ID == "" {
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

func ConfigureCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
}
