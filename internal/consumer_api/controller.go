package consumer_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/scenario"
	"switchboard/internal/workspace_setting"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type requestDetails struct {
	Method   string              `json:"method"`
	Hostname string              `json:"hostname"`
	Port     string              `json:"port"`
	Host     string              `json:"host"`
	Path     string              `json:"path"`
	Body     string              `json:"body"`
	Headers  map[string][]string `json:"headers"`
}

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

func getRequestDetails(req *http.Request) (*requestDetails, error) {
	bodyBytes, errBody := ioutil.ReadAll(req.Body)
	if errBody != nil {
		return nil, errBody
	}
	return &requestDetails{
		Method:   req.Method,
		Port:     req.URL.Port(),
		Hostname: req.URL.Hostname(),
		Host:     req.Host,
		Path:     req.URL.Path,
		Headers:  req.Header,
		Body:     string(bodyBytes),
	}, nil
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
				scenarioFound := false
				for _, sc := range ec.ScenarioConfigs {
					if scenarioFound {
						break
					}
					if sc.IsActive {
						activeScenarioId = sc.ScenarioID
						scenarioFound = true
					}
					// Run interception rules
					for _, ir := range ec.InterceptionRules {
						req, errR := getRequestDetails(c.Request)
						reqData, errM := json.Marshal(req)
						if errR != nil || errM != nil {
							logrus.Errorln(errM, errR)
							c.JSON(http.StatusInternalServerError, gin.H{
								"error": "could not read request details",
							})
							return
						}
						matched := common.ApplyJsonLogic(ir.MatcherExpression, string(reqData))
						if matched {
							activeScenarioId = ir.TargetScenarioId
							scenarioFound = true
						}
					}
				}
				break
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
