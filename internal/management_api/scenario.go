package management_api

import (
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateScenarioRoute(c *gin.Context) {
	var payload models.CreateScenarioRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, err_utils.NewDetailedError(
			err_utils.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	createdScenario, createErr := db.CreateScenario(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdScenario)
		return
	}

	if createErr.ErrorCode == err_utils.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "duplicate scenario"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}
