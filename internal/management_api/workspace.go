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

func CreateWorkspaceRoute(c *gin.Context) {
	var payload models.CreateWorkspaceRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, common.NewDetailedError(
			common.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	createdWs, createErr := db.CreateWorkspace(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdWs)
		return
	}

	if createErr.ErrorCode == common.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "workspace already exists"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func GetWorkspacesRoute(c *gin.Context) {
	ws, err := db.GetWorkspaces()
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve workspaces. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, ws)
}

func GetUserWorkspacesRoute(c *gin.Context) {
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	ws, err := db.GetUserWorkspaces(currentUser.ID)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve user workspaces. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, ws)
}

func DeleteWorkspaceRoute(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	currentUser := c.Value(common.REQ_USER_KEY).(*models.User)
	ok, err := db.DeleteWorkspace(currentUser.ID, workspaceID)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not delete workspace %s. [error code: %s] [description: %s]", workspaceID, err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	if !ok {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
