package models

import "time"

type HTTPResponseScenarioConfig struct {
	StatusCode              uint16 `json:"statusCode" bson:"statusCode,omitempty"`
	ResponseBodyTemplate    string `json:"responseBodyTemplate" bson:"responseBodyTemplate,omitempty"`
	ResponseHeadersTemplate string `json:"responseHeadersTemplate" bson:"responseHeadersTemplate,omitempty"`
}

type ProxyScenarioConfig struct {
	UpstreamID    string            `json:"upstreamID" bson:"upstreamID,omitempty"`
	InjectHeaders map[string]string `json:"injectHeaders" bson:"injectHeaders,omitempty"`
}

type NetworkScenarioConfig struct {
	Type string `json:"type" bson:"type,omitempty"`
}

type Scenario struct {
	ID                         string                      `json:"id" bson:"id,omitempty"`
	EndpointId                 string                      `json:"endpointId" bson:"endpointId,omitempty"`
	Type                       string                      `json:"type" bson:"type,omitempty"`
	IsDefault                  bool                        `json:"isDefault" bson:"isDefault"`
	HTTPResponseScenarioConfig *HTTPResponseScenarioConfig `json:"httpResponseScenarioConfig,omitempty" bson:"httpResponseScenarioConfig,omitempty"`
	ProxyScenarioConfig        *ProxyScenarioConfig        `json:"proxyScenarioConfig,omitempty" bson:"proxyScenarioConfig,omitempty"`
	NetworkScenarioConfig      *NetworkScenarioConfig      `json:"networkScenarioConfig,omitempty" bson:"networkScenarioConfig,omitempty"`
	CreatedBy                  string                      `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt                  time.Time                   `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt                  time.Time                   `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type HTTPResponseScenarioRequestBody struct {
	StatusCode              uint16 `json:"statusCode"`
	ResponseBodyTemplate    string `json:"responseBodyTemplate"`
	ResponseHeadersTemplate string `json:"responseHeadersTemplate"`
}

type ProxyScenarioRequestBody struct {
	UpstreamID    string            `json:"upstreamID"`
	InjectHeaders map[string]string `json:"injectHeaders"`
}

type NetworkScenarioRequestBody struct {
	Type string `json:"type"`
}

type CreateScenarioRequestBody struct {
	EndpointId                 string                           `json:"endpointId" binding:"required"`
	Type                       string                           `json:"type" binding:"required,validScenario"`
	HTTPResponseScenarioConfig *HTTPResponseScenarioRequestBody `json:"httpResponseScenarioConfig,omitempty"`
	ProxyScenarioConfig        *ProxyScenarioRequestBody        `json:"proxyScenarioConfig,omitempty"`
	NetworkScenarioConfig      *NetworkScenarioRequestBody      `json:"networkScenarioConfig,omitempty"`
}
