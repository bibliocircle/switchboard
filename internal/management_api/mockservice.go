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

func CreateMockServiceRoute(c *gin.Context) {
	var payload models.CreateMockServiceRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, common.NewDetailedError(
			common.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	createdMockService, createErr := db.CreateMockService(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdMockService)
		return
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "mock service already exists"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func DeleteMockServiceRoute(c *gin.Context) {
	mockServiceID := c.Param("mockServiceId")
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	ok, err := db.DeleteMockService(currentUser.ID, mockServiceID)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not delete mock service %s. [error code: %s] [description: %s]", mockServiceID, err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	if !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
