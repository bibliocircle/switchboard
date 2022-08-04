package consumer_api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/scenario"
	"switchboard/internal/workspace_setting"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func activateHTTPResponseScenario(cfg *scenario.HTTPResponseScenarioConfig, ctx *gin.Context) {
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

func activateScenario(sc *scenario.Scenario, ctx *gin.Context) {
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

		wss, wsErr := workspace_setting.GetWorkspaceSetting(wsID, msID)
		if wsErr != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if wss == nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": fmt.Sprintf("mock service '%s' is not enabled on this workspace", msID),
			})
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

		sc, scErr := scenario.GetScenarioByID(activeScenarioId)
		if scErr != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if sc == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "active scenario could not be found for the endpoint on this workspace",
			})
			return
		}

		activateScenario(sc, c)
	}
}
