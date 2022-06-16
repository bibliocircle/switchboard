package upstream

import (
	"encoding/json"
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/db"
	"switchboard/internal/common/err_utils"

	"github.com/gin-gonic/gin"
)

func CreateUpstreamRoute(c *gin.Context) {
	var payload Upstream
	decodeErr := json.NewDecoder(c.Request.Body).Decode(&payload)
	if decodeErr != nil {
		c.JSON(http.StatusBadRequest, err_utils.NewDetailedError(
			err_utils.ErrorUnparsablePayload,
			"could not parse request payload",
		))
		return
	}
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	createdUpstream, createErr := CreateUpstream(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdUpstream)
		return
	}
	wrappedErr := db.GetDbError(createErr)

	if wrappedErr.ErrorCode == err_utils.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, "endpoint already exists")
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}
