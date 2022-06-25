package management_api

import (
	"fmt"
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func CreateEndpointRoute(c *gin.Context) {
	var payload models.CreateEndpointRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, common.NewDetailedError(
			common.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	createdEndpoint, createErr := db.CreateEndpoint(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdEndpoint)
		return
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "endpoint already exists"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func GetEndpointsByMockServiceIdRoute(c *gin.Context) {
	endpoints, err := db.GetEndpoints(c.Param("mockServiceId"))
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve endpoints. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, endpoints)
}

func DeleteEndpointRoute(c *gin.Context) {
	endpointID := c.Param("endpointId")
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	ok, err := db.DeleteEndpoint(currentUser.ID, endpointID)
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
