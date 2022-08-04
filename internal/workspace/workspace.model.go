package workspace

import "time"

type Workspace struct {
	ID        string     `json:"id" bson:"id,omitempty"`
	Name      string     `json:"name" bson:"name,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	CreatedBy string     `json:"createdBy" bson:"createdBy,omitempty"`
	CreatedAt *time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type CreateWorkspaceRequestBody struct {
	Name      string `json:"name" binding:"required"`
	ExpiresAt string `json:"expiresAt,omitempty" binding:"omitempty,isodate"`
}
