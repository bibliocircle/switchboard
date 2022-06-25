package management_api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type SignUpPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginRoute(c *gin.Context) {
	var payload LoginPayload
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(http.StatusBadRequest, common.DetailedError{
			ErrorCode:   common.ErrorUnparsablePayload,
			Description: "Request body could not be parsed",
		})
	}
	if payload.Email == "" || payload.Password == "" {
		c.JSON(http.StatusBadRequest, common.DetailedError{
			ErrorCode:   common.ErrorPayloadMissingRequiredFields,
			Description: "One or more of the required fields are missing",
		})
		return
	}
	userEntity, err := db.GetUserByEmailPassword(payload.Email, payload.Password)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userEntity == nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, tokenError := common.CreateSignedAuthToken(*userEntity)
	if tokenError != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenExpiry, parseErr := strconv.ParseInt(os.Getenv("AUTH_TOKEN_EXPIRY_SECONDS"), 10, 32)
	if parseErr != nil {
		logrus.Error("could not parse AUTH_TOKEN_EXPIRY_SECONDS as int32")
		tokenExpiry = 86400 // default
	}
	c.SetCookie(os.Getenv("AUTH_COOKIE_NAME"), *token, int(tokenExpiry), "/", os.Getenv("AUTH_COOKIE_DOMAIN"), false, false)
	c.JSON(http.StatusOK, LoginResponse{
		Token: *token,
	})
}

func SignUpRoute(c *gin.Context) {
	var payload SignUpPayload
	if err := json.NewDecoder(c.Request.Body).Decode(&payload); err != nil {
		c.JSON(http.StatusBadRequest, common.DetailedError{
			ErrorCode:   common.ErrorUnparsablePayload,
			Description: "Request body could not be parsed",
		})
	}
	if payload.Email == "" || payload.Password == "" {
		c.JSON(http.StatusBadRequest, common.DetailedError{
			ErrorCode:   common.ErrorPayloadMissingRequiredFields,
			Description: "One or more of the required fields are missing",
		})
	}

	createdUser, err := db.CreateUser(&models.CreateUserRequest{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	})
	if err != nil {
		if err.ErrorCode == common.ErrorDuplicateEntity {
			c.JSON(http.StatusConflict, common.DetailedError{
				ErrorCode:   common.ErrorUserExists,
				Description: "A matching user already exists",
			})
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, models.User{
		ID:        createdUser.ID,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	})
}

func LogOutRoute(c *gin.Context) {
	c.SetCookie(os.Getenv("AUTH_COOKIE_NAME"), "", -1, "/", os.Getenv("AUTH_COOKIE_DOMAIN"), false, false)
}

func ResetPassword(c *gin.Context) {
	// TODO
}
