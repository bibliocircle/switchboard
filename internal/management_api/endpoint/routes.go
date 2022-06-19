package endpoint

import (
	"fmt"
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateEndpointRequestBody struct {
	MockServiceId string `json:"mockServiceId" binding:"required"`
	Path          string `json:"path" binding:"required,absolutePath"`
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

func GetEndpointsByMockServiceIdRoute(c *gin.Context) {
	endpoints, err := GetEndpoints(c.Param("mockServiceId"))
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve endpoints. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, endpoints)
}

func DeleteEndpointRoute(c *gin.Context) {
	endpointID := c.Param("endpointId")
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	ok, err := DeleteEndpoint(currentUser.ID, endpointID)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not delete endpoint %s. [error code: %s] [description: %s]", endpointID, err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	if !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
