package management_api

import (
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateScenarioRoute(c *gin.Context) {
	var payload models.CreateScenarioRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, common.NewDetailedError(
			common.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	createdScenario, createErr := db.CreateScenario(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdScenario)
		return
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "duplicate scenario"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}
