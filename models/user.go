package models

type User struct {
	UUID     string  `json:"uuid" bson:"_id, omitempty"`
	Username string  `json:"username" bson:"username, omitempty"`
	ImageUrl string  `json:"image_url" bson:"image_url"`
	Email    string  `json:"email" bson:"email, omitempty"`
	Password string  `bson:"password, omitempty"`
	Goal     float64 `json:"goal" bson:"goal, omitempty"`
	UserRole Role    `json:"role" bson:"role, omitempty"`
	Runs     []Run   `json:"runs" bson:"runs, omitempty"`
}

type ReturnUser struct {
	UUID     string  `json:"uuid" bson:"_id, omitempty"`
	Email    string  `json:"email" bson:"email, omitempty"`
	Username string  `json:"username" bson:"username, omitempty"`
	Goal     float64 `json:"goal" bson:"goal, omitempty"`
	ImageUrl string  `json:"image_url" bson:"image_url"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleNone   Role = "none"
)

type LeaderboardUser struct {
	UUID     string  `json:"uuid" bson:"_id, omitempty"`
	Email    string  `json:"email" bson:"email, omitempty"`
	Username string  `json:"username" bson:"username, omitempty"`
	Total    float64 `json:"total" bson:"total, omitempty"`
	ImageUrl string  `json:"image_url" bson:"image_url"`
}
