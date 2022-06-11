package scenario

import "time"

type ScenarioConfig struct {
	ID                   string            `json:"id" bson:"id,omitempty"`
	StatusCode           int32             `json:"statusCode" bson:"statusCode,omitempty"`
	ResponseBodyTemplate string            `json:"responseBodyTemplate" bson:"responseBodyTemplate,omitempty"`
	ResponseHeaders      map[string]string `json:"responseHeaders" bson:"responseHeaders,omitempty"`
}

type Scenario struct {
	ID        string         `json:"id" bson:"id,omitempty"`
	Type      string         `json:"type" bson:"type,omitempty"`
	Config    ScenarioConfig `json:"config" bson:"config,omitempty"`
	CreatedAt time.Time      `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt" bson:"updatedAt,omitempty"`
}
