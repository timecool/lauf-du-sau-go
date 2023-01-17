package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id, omitempty"`
	Username string             `json:"username,omitempty" bson:"username, omitempty"`
	ImageUrl string             `json:"image_url" bson:"image_url"`
	Email    string             `json:"email,omitempty" bson:"email, omitempty"`
	Password string             `json:"password" bson:"password, omitempty" `
	Goal     float64            `json:"goal" bson:"goal, omitempty"`
	UserRole Role               `json:"role,omitempty" bson:"role, omitempty"`
}

type ReturnUser struct {
	ID       primitive.ObjectID `json:"id" bson:"_id, omitempty"`
	Email    string             `json:"email" bson:"email, omitempty"`
	Username string             `json:"username" bson:"username, omitempty"`
	Goal     float64            `json:"goal" bson:"goal, omitempty"`
	ImageUrl string             `json:"image_url" bson:"image_url"`
	UserRole Role               `json:"role" bson:"role, omitempty"`
}

type ResetPassword struct {
	Password string `bson:"password, omitempty"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleNone   Role = "none"
)

type LeaderboardUserID struct {
	UserID string  `json:"user_id" bson:"_id, omitempty"`
	Total  float64 `json:"total" bson:"total, omitempty"`
}

type LeaderboardUser struct {
	User  ReturnUser `json:"user" bson:"user, omitempty"`
	Total float64    `json:"total" bson:"total, omitempty"`
}
