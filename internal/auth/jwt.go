package auth

import (
	"os"
	"strconv"
	"switchboard/internal/user"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

func CreateSignedAuthToken(user user.User) (*string, error) {
	type AuthClaims struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		jwt.StandardClaims
	}

	tokenExpiry, err := strconv.ParseInt(os.Getenv("AUTH_TOKEN_EXPIRY_SECONDS"), 10, 32)
	if err != nil {
		log.Info("could not parse AUTH_TOKEN_EXPIRY_SECONDS environment variable to int64")
		tokenExpiry = 86400 // default
	}

	claims := AuthClaims{
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + tokenExpiry,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey := []byte(os.Getenv("AUTH_TOKEN_KEY"))
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}
