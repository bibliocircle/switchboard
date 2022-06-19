package workspace

import (
	"fmt"
	"net/http"
	"switchboard/internal/common/auth"
	"switchboard/internal/common/constants"
	"switchboard/internal/common/err_utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CreateWorkspaceRequestBody struct {
	Name      string `json:"name" binding:"required"`
	ExpiresAt string `json:"expiresAt,omitempty" binding:"omitempty,isodate"`
}

func CreateWorkspaceRoute(c *gin.Context) {
	var payload CreateWorkspaceRequestBody
	if bindErr := c.ShouldBindJSON(&payload); bindErr != nil {
		c.JSON(http.StatusBadRequest, err_utils.NewDetailedError(
			err_utils.ErrorUnparsablePayload,
			bindErr.Error(),
		))
		return
	}
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	createdWs, createErr := CreateWorkspace(currentUser.ID, &payload)
	if createErr == nil {
		c.JSON(http.StatusCreated, createdWs)
		return
	}

	if createErr.ErrorCode == err_utils.ErrorDuplicateEntity {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "workspace already exists"})
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func GetWorkspacesRoute(c *gin.Context) {
	ws, err := GetWorkspaces()
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve workspaces. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, ws)
}

func GetUserWorkspacesRoute(c *gin.Context) {
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	ws, err := GetUserWorkspaces(currentUser.ID)
	if err != nil {
		log.Errorln(fmt.Sprintf("could not retrieve user workspaces. [error code: %s] [description: %s]", err.ErrorCode, err.Description))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.ErrorCode})
		return
	}
	c.JSON(http.StatusOK, ws)
}

func DeleteWorkspaceRoute(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	currentUser := c.Value(constants.REQ_USER_KEY).(*auth.User)
	ok, err := DeleteWorkspace(currentUser.ID, workspaceID)
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
