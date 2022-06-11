package endpoint

import (
	"encoding/json"
	"net/http"
	"switchboard/internal/common"

	"github.com/gin-gonic/gin"
)

func CreateEndpointRoute(c *gin.Context) {
	var payload Endpoint
	decodeErr := json.NewDecoder(c.Request.Body).Decode(&payload)
	if decodeErr != nil {
		c.JSON(http.StatusBadRequest, common.NewDetailedError(
			common.ErrorUnparsablePayload,
			"could not parse request payload",
		))
		return
	}

}
