package mockservice

import (
	"switchboard/internal/scenario"
	"time"
)

type GlobalMockServiceConfig struct {
	InjectHeaders       []scenario.HTTPHeader `json:"injectHeaders" bson:"injectHeaders"`
	GlobalResponseDelay uint16                `json:"globalResponseDelay" bson:"responseDelay"`
}

type MockService struct {
	ID        string                  `json:"id" bson:"id,omitempty"`
	Name      string                  `json:"name" bson:"name,omitempty"`
	Type      string                  `json:"type" bson:"type,omitempty"`
	Config    GlobalMockServiceConfig `json:"config" bson:"config,omitempty"`
	CreatedBy string                  `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt *time.Time              `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt *time.Time              `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type CreateMockServiceRequestBody struct {
	ID     string                  `json:"id" binding:"required"`
	Name   string                  `json:"name" binding:"required"`
	Type   string                  `json:"type" binding:"required"`
	Config GlobalMockServiceConfig `json:"config" binding:"required"`
}
