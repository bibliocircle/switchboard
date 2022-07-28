package consumer_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/db"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func activateHTTPResponseScenario(cfg *models.HTTPResponseScenarioConfig, ctx *gin.Context) {
	ctx.Status(int(cfg.StatusCode))
	headers := map[string]string{}
	headersStr, rhErr := common.GenFakeData(cfg.ResponseHeadersTemplate)
	if rhErr != nil {
		log.Errorln(rhErr)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	bodyStr, bErr := common.GenFakeData(cfg.ResponseBodyTemplate)
	if bErr != nil {
		log.Errorln(bErr)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.Unmarshal([]byte(headersStr), &headers)
	for h_name, h_value := range headers {
		ctx.Header(h_name, h_value)
	}
	ctx.Writer.Write([]byte(bodyStr))
}

func activateScenario(sc *models.Scenario, ctx *gin.Context) {
	switch sc.Type {
	case common.HTTP_SCENARIO_TYPE:
		activateHTTPResponseScenario(sc.HTTPResponseScenarioConfig, ctx)
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("misconfigured service! unknown scenario type: %s", sc.Type),
		})
	}
}

func CreateRoute(mockServiceID, endpointID string) gin.HandlerFunc {

	return func(c *gin.Context) {
		wsID := c.Param("workspaceId")
		msID := mockServiceID
		eID := endpointID

		wss, wsErr := db.GetWorkspaceSetting(wsID, msID)
		if wsErr != nil {
			if wsErr.Error() == mongo.ErrNoDocuments.Error() {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": fmt.Sprintf("mock service '%s' is not enabled on this workspace", msID),
				})
				return
			}
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		var activeScenarioId string
		for _, ec := range wss.EndpointConfigs {
			if ec.EndpointID == eID {
				for _, sc := range ec.ScenarioConfigs {
					if sc.IsActive {
						activeScenarioId = sc.ScenarioID
						break
					}
				}
			}
		}

		if activeScenarioId == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "no active scenario configured for this endpoint on this workspace",
			})
			return
		}

		sc, scErr := db.GetScenarioByID(activeScenarioId)
		if scErr != nil {
			if scErr.Error() == mongo.ErrNoDocuments.Error() {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "active scenario could not be found for the endpoint on this workspace",
				})
				return
			}
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		activateScenario(sc, c)
	}
}
