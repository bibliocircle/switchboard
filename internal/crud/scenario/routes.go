package scenario

import (
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"

	"github.com/gin-gonic/gin"
)

type CreateScenarioConfigRequestBody struct {
	StatusCode           int32             `json:"statusCode" binding:"required"`
	ResponseBodyTemplate string            `json:"responseBodyTemplate"`
	ResponseHeaders      map[string]string `json:"responseHeaders"`
}

type CreateScenarioRequestBody struct {
	EndpointId string                          `json:"endpointId" binding:"required"`
	Type       string                          `json:"type" binding:"required"`
	Config     CreateScenarioConfigRequestBody `json:"config" binding:"required"`
}

func CreateScenarioRoute(c *gin.Context) {
	var payload CreateScenarioRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, err_utils.NewDetailedError(
			err_utils.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	createdScenario, createErr := CreateScenario(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdScenario)
		return
	}

	if createErr.ErrorCode == err_utils.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "duplicate scenario"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}
