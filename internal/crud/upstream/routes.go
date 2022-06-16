package upstream

import (
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusUnprocessableEntity, "endpoint already exists")
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}
