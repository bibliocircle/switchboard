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

func CreateUpstreamRoute(c *gin.Context) {
	var payload models.CreateUpstreamRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, common.NewDetailedError(
			common.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	createdUpstream, createErr := db.CreateUpstream(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdUpstream)
		return
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "endpoint already exists"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func DeleteUpstreamRoute(c *gin.Context) {
	upstreamID := c.Param("upstreamId")
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	ok, err := db.DeleteUpstream(currentUser.ID, upstreamID)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not delete upstream %s. [error code: %s] [description: %s]", upstreamID, err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	if !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
