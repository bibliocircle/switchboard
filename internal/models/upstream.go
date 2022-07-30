package models

import "time"

type Upstream struct {
	ID            string     `json:"id" bson:"id,omitempty"`
	MockServiceId string     `json:"mockServiceId" bson:"mockServiceId,omitempty"`
	Name          string     `json:"name" bson:"name,omitempty"`
	URL           string     `json:"url" bson:"url,omitempty"`
	CreatedBy     string     `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt     *time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     *time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type CreateUpstreamRequestBody struct {
	MockServiceId string `json:"mockServiceId" binding:"required"`
	Name          string `json:"name" binding:"required"`
	URL           string `json:"url" binding:"required,url"`
}
