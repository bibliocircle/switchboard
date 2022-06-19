package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/management_api/user"

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
		c.JSON(http.StatusBadRequest, err_utils.DetailedError{
			ErrorCode:   err_utils.ErrorUnparsablePayload,
			Description: "Request body could not be parsed",
		})
	}
	if payload.Email == "" || payload.Password == "" {
		c.JSON(http.StatusBadRequest, err_utils.DetailedError{
			ErrorCode:   err_utils.ErrorPayloadMissingRequiredFields,
			Description: "One or more of the required fields are missing",
		})
		return
	}
	userEntity, err := user.GetUserByEmailPassword(payload.Email, payload.Password)
	if err != nil {
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userEntity == nil {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, tokenError := CreateSignedAuthToken(*userEntity)
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
		c.JSON(http.StatusBadRequest, err_utils.DetailedError{
			ErrorCode:   err_utils.ErrorUnparsablePayload,
			Description: "Request body could not be parsed",
		})
	}
	if payload.Email == "" || payload.Password == "" {
		c.JSON(http.StatusBadRequest, err_utils.DetailedError{
			ErrorCode:   err_utils.ErrorPayloadMissingRequiredFields,
			Description: "One or more of the required fields are missing",
		})
	}

	createdUser, err := user.CreateUser(&user.CreateUserRequest{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  payload.Password,
	})
	if err != nil {
		if err.ErrorCode == err_utils.ErrorDuplicateEntity {
			c.JSON(http.StatusConflict, err_utils.DetailedError{
				ErrorCode:   err_utils.ErrorUserExists,
				Description: "A matching user already exists",
			})
			return
		}
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, user.User{
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
