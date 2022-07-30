package models

import "time"

type Endpoint struct {
	ID            string     `json:"id" bson:"id,omitempty"`
	MockServiceId string     `json:"mockServiceId" bson:"mockServiceId,omitempty"`
	Path          string     `json:"path" bson:"path,omitempty"`
	Method        string     `json:"method" bson:"method,omitempty"`
	Description   string     `json:"description" bson:"description,omitempty"`
	ResponseDelay uint16     `json:"responseDelay" bson:"responseDelay,omitempty"`
	CreatedBy     string     `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt     *time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     *time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type CreateEndpointRequestBody struct {
	MockServiceId string `json:"mockServiceId" binding:"required"`
	Path          string `json:"path" binding:"required,absolutePath"`
	Method        string `json:"method" binding:"required"`
	Description   string `json:"description" binding:"required"`
	ResponseDelay uint16 `json:"responseDelay"`
}
