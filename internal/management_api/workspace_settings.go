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

func AddMockServiceToWorkspaceRoute(c *gin.Context) {
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	workspaceId := c.Param("workspaceId")
	mockServiceId := c.Param("mockServiceId")

	isWsOwner, errPerm := db.IsWorkspaceOwner(currentUser.ID, workspaceId)
	if errPerm != nil {
		log.Errorln(fmt.Sprintf("could not check workspace ownership. [error code: %s] [description: %s]", errPerm.ErrorCode, errPerm.Description))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isWsOwner {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}

	err := db.AddMockServiceToWorkspace(currentUser.ID, workspaceId, mockServiceId)
	if err == nil {
		c.Writer.WriteHeader(http.StatusCreated)
		return
	}

	if err.ErrorCode == common.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "mock service is already added to the workspace"})
		return
	}

	log.Errorln(fmt.Sprintf("could not add mock service to workspace. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func ActivateMockServiceScenarioRoute(c *gin.Context) {
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	workspaceId := c.Param("workspaceId")
	mockServiceId := c.Param("mockServiceId")
	endpointId := c.Param("endpointId")
	scenarioId := c.Param("scenarioId")

	isWsOwner, errPerm := db.IsWorkspaceOwner(currentUser.ID, workspaceId)
	if errPerm != nil {
		log.Errorln(fmt.Sprintf("could not check workspace ownership. [error code: %s] [description: %s]", errPerm.ErrorCode, errPerm.Description))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isWsOwner {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}

	ok, updateErr := db.ActivateMockServiceScenario(workspaceId, mockServiceId, endpointId, scenarioId)
	if updateErr != nil {
		log.Errorln(fmt.Sprintf("could not activate scenario in workspace. [error code: %s] [description: %s]", updateErr.ErrorCode, updateErr.Description))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}

func UpdateWsMockServiceConfigRoute(c *gin.Context) {
	var payload models.UpdateMockServiceConfigRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, common.NewDetailedError(
			common.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	workspaceId := c.Param("workspaceId")
	mockServiceId := c.Param("mockServiceId")
	endpointId := c.Param("endpointId")

	isWsOwner, errPerm := db.IsWorkspaceOwner(currentUser.ID, workspaceId)
	if errPerm != nil {
		log.Errorln(fmt.Sprintf("could not check workspace ownership. [error code: %s] [description: %s]", errPerm.ErrorCode, errPerm.Description))
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isWsOwner {
		c.Writer.WriteHeader(http.StatusForbidden)
		return
	}

	updatedWss, updateErr := db.UpdateWsMockServiceConfig(workspaceId, mockServiceId, endpointId, &payload)
	if updateErr == nil {
		c.JSON(http.StatusOK, updatedWss)
		return
	}

	log.Errorln(fmt.Sprintf("could not update workspace setting. [error code: %s] [description: %s]", updateErr.ErrorCode, updateErr.Description))
	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func GetWorkspaceSettingsRoute(c *gin.Context) {
	ws, err := db.GetWorkspaceSettings(c.Param("workspaceId"))
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve workspace settings. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, ws)
}
