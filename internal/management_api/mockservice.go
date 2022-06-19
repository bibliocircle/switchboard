package management_api

import (
	"fmt"
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func CreateMockServiceRoute(c *gin.Context) {
	var payload models.CreateMockServiceRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, err_utils.NewDetailedError(
			err_utils.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	createdMockService, createErr := db.CreateMockService(currentUser.ID, &payload)
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

func GetMockServicesRoute(c *gin.Context) {
	mockServices, err := db.GetMockServices()
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve mock services. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, mockServices)
}

func DeleteMockServiceRoute(c *gin.Context) {
	mockServiceID := c.Param("mockServiceId")
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
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
