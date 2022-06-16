package endpoint

import (
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"

	"github.com/gin-gonic/gin"
)

type CreateEndpointRequestBody struct {
	MockServiceId string `json:"mockServiceId" binding:"required"`
	Path          string `json:"path" binding:"required"`
	Method        string `json:"method" binding:"required"`
	Description   string `json:"description" binding:"required"`
	ResponseDelay int64  `json:"responseDelay"`
}

func CreateEndpointRoute(c *gin.Context) {
	var payload CreateEndpointRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, err_utils.NewDetailedError(
			err_utils.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	createdEndpoint, createErr := CreateEndpoint(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdEndpoint)
		return
	}

	if createErr.ErrorCode == err_utils.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "endpoint already exists"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}
