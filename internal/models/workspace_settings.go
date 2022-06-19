package models

import "time"

type ScenarioConfig struct {
	ScenarioID string `json:"scenarioId" bson:"id,omitempty"`
	IsActive   bool   `json:"isActive" bson:"isActive"`
}

type WorkspaceSetting struct {
	WorkspaceID   string           `json:"workspaceId" bson:"workspaceId,omitempty"`
	MockServiceID string           `json:"mockServiceId" bson:"mockServiceId,omitempty"`
	EndpointID    string           `json:"endpointId" bson:"endpointId,omitempty"`
	Scenarios     []ScenarioConfig `json:"scenarios" bson:"scenarios"`
	ResponseDelay uint16           `json:"responseDelay" bson:"responseDelay,omitempty"`
	CreatedBy     string           `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt     time.Time        `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time        `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type UpdateMockServiceConfigRequestBody struct {
	ResponseDelay uint16 `json:"responseDelay" binding:"required"`
}
