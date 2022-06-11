package endpoint

import "time"

type Endpoint struct {
	ID            string    `json:"id" bson:"id,omitempty"`
	Path          string    `json:"path" bson:"path,omitempty"`
	Method        string    `json:"method" bson:"method,omitempty"`
	ResponseDelay int64     `json:"responseDelay" bson:"responseDelay,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}
