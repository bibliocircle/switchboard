package models

type ScenarioConfig struct {
	ScenarioID string `json:"scenarioId" bson:"id,omitempty"`
	IsActive   bool   `json:"isActive" bson:"isActive"`
}

type EndpointConfig struct {
	EndpointID      string           `json:"endpointId" bson:"endpointId,omitempty"`
	ScenarioConfigs []ScenarioConfig `json:"scenarioConfigs" bson:"scenarioConfigs"`
	ResponseDelay   uint16           `json:"responseDelay" bson:"responseDelay,omitempty"`
}

type WorkspaceSetting struct {
	WorkspaceID     string                  `json:"workspaceId" bson:"workspaceId,omitempty"`
	MockServiceID   string                  `json:"mockServiceId" bson:"mockServiceId,omitempty"`
	Config          GlobalMockServiceConfig `json:"config" bson:"config,omitempty"`
	EndpointConfigs []EndpointConfig        `json:"endpointConfigs" bson:"endpointConfigs"`
}

type UpdateMockServiceConfigRequestBody struct {
	ResponseDelay uint16 `json:"responseDelay" binding:"required"`
}
