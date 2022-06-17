package upstream

import (
	"fmt"
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateUpstreamRequestBody struct {
	MockServiceId string `json:"mockServiceId" binding:"required"`
	Name          string `json:"name" binding:"required"`
	URL           string `json:"url" binding:"required,url"`
}

func CreateUpstreamRoute(c *gin.Context) {
	var payload CreateUpstreamRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, err_utils.NewDetailedError(
			err_utils.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	createdUpstream, createErr := CreateUpstream(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdUpstream)
		return
	}

	if createErr.ErrorCode == err_utils.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "endpoint already exists"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func GetUpstreamsByMockServiceIdRoute(c *gin.Context) {
	upstreams, err := GetUpstreams(c.Param("mockServiceId"))
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve upstreams. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, upstreams)
}

func DeleteUpstreamRoute(c *gin.Context) {
	upstreamID := c.Param("upstreamId")
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	ok, err := DeleteUpstream(currentUser.ID, upstreamID)
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
