package workspace_setting

import "switchboard/internal/mockservice"

type InterceptionRule struct {
	ID                string `json:"id" bson:"id,omitempty"`
	Name              string `json:"name" bson:"name,omitempty"`
	MatcherExpression string `json:"matcherExpression" bson:"matcherExpression"`
	TargetScenarioId  string `json:"targetScenarioId" bson:"expression,targetScenarioId"`
}

type ScenarioConfig struct {
	ID         string `json:"id" bson:"id"`
	ScenarioID string `json:"scenarioId" bson:"scenarioId,omitempty"`
	IsActive   bool   `json:"isActive" bson:"isActive"`
}

type EndpointConfig struct {
	ID                string             `json:"id" bson:"id"`
	EndpointID        string             `json:"endpointId" bson:"endpointId,omitempty"`
	ScenarioConfigs   []ScenarioConfig   `json:"scenarioConfigs" bson:"scenarioConfigs"`
	InterceptionRules []InterceptionRule `json:"interceptionRules" bson:"interceptionRules"`
	ResponseDelay     uint16             `json:"responseDelay" bson:"responseDelay,omitempty"`
}

type WorkspaceSetting struct {
	ID              string                              `json:"id" bson:"id"`
	WorkspaceID     string                              `json:"workspaceId" bson:"workspaceId,omitempty"`
	MockServiceID   string                              `json:"mockServiceId" bson:"mockServiceId,omitempty"`
	Config          mockservice.GlobalMockServiceConfig `json:"config" bson:"config,omitempty"`
	EndpointConfigs []EndpointConfig                    `json:"endpointConfigs" bson:"endpointConfigs"`
}

type CreateInterceptionRuleRequestBody struct {
	Name              string `json:"name" bson:"name,omitempty"`
	MatcherExpression string `json:"matcherExpression" bson:"matcherExpression"`
	TargetScenarioId  string `json:"targetScenarioId" bson:"expression,targetScenarioId"`
}

type UpdateMockServiceConfigRequestBody struct {
	ResponseDelay uint16 `json:"responseDelay" binding:"required"`
}
