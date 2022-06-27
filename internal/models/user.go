package models

import "time"

type User struct {
	ID            string    `json:"id" bson:"id,omitempty" mapstructure:"id"`
	FirstName     string    `json:"firstName" bson:"firstName,omitempty" mapstructure:"firstName"`
	LastName      string    `json:"lastName" bson:"lastName,omitempty" mapstructure:"lastName"`
	Email         string    `json:"email" bson:"email,omitempty" mapstructure:"email"`
	Password      string    `json:"-" bson:"password,omitempty"`
	GoogleUserId  string    `json:"googleUserId" bson:"googleUserId,omitempty"`
	EmailVerified string    `json:"emailVerified" bson:"emailVerified,omitempty"`
	Deleted       string    `json:"-" bson:"deleted,omitempty"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}
