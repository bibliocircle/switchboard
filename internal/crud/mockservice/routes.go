package mockservice

import (
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"

	"github.com/gin-gonic/gin"
)

type CreateMockServiceRequestBody struct {
	Name   string                  `json:"name" binding:"required"`
	Key    string                  `json:"key" binding:"required"`
	Type   string                  `json:"type" binding:"required"`
	Config GlobalMockServiceConfig `json:"config" binding:"required"`
}

func CreateMockServiceRoute(c *gin.Context) {
	var payload CreateMockServiceRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, err_utils.NewDetailedError(
			err_utils.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	createdMockService, createErr := CreateMockService(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdMockService)
		return
	}

	if createErr.ErrorCode == err_utils.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "mock service already exists"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}
