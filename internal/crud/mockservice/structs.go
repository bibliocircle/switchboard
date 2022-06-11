package mockservice

import (
	"switchboard/internal/crud/endpoint"
	"switchboard/internal/crud/upstream"
	"time"
)

type CorsConfig struct {
	AllowedOrigins []string `json:"allowedOrigins" bson:"allowedOrigins,omitempty"`
	AllowedMethods []string `json:"allowedMethods" bson:"allowedMethods,omitempty"`
	AllowedHeaders []string `json:"allowedHeaders" bson:"allowedHeaders,omitempty"`
}

type GlobalMockServiceConfig struct {
	CORS          CorsConfig          `json:"cors" bson:"cors,omitempty"`
	InjectHeaders map[string]string   `json:"injectHeaders" bson:"injectHeaders,omitempty"`
	Upstreams     []upstream.Upstream `json:"upstreams" bson:"upstreams,omitempty"`
}

type MockService struct {
	ID        string                  `json:"id" bson:"id,omitempty"`
	Title     string                  `json:"title" bson:"title,omitempty"`
	Version   string                  `json:"version" bson:"version,omitempty"`
	Type      string                  `json:"type" bson:"type,omitempty"`
	Endpoints []endpoint.Endpoint     `json:"endpoints" bson:"endpoints,omitempty"`
	Config    GlobalMockServiceConfig `json:"config" bson:"config,omitempty"`
	CreatedAt time.Time               `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time               `json:"updatedAt" bson:"updatedAt,omitempty"`
}
