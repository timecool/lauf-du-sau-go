package models

type User struct {
	UUID     string  `json:"uuid" bson:"_id, omitempty"`
	Email    string  `json:"email" bson:"email, omitempty"`
	Password string  `bson:"password, omitempty"`
	Goal     float64 `json:"goal" bson:"goal, omitempty"`
	UserRole Role    `json:"role" bson:"role, omitempty"`
	Runs     []Run   `json:"runs" bson:"runs"`
}

type ReturnUser struct {
	UUID  string `json:"uuid" bson:"_id, omitempty"`
	Email string `json:"email" bson:"email, omitempty"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleNone   Role = "none"
)
