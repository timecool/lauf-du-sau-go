package models

import (
	"time"
)

type Run struct {
	UUID     string       `json:"uuid" bson:"_id, omitempty"`
	Date     time.Time    `json:"date" bson:"date, omitempty"`
	CreateAt time.Time    `json:"create_at" bson:"create_at, omitempty"`
	Time     float64      `json:"time" bson:"time, omitempty"`
	Distance float64      `json:"distance" bson:"distance, omitempty"`
	Url      string       `json:"url" bson:"url, omitempty"`
	Status   RunStatus    `json:"status" bson:"status, omitempty"`
	Messages []RunMessage `json:"messages" bson:"messages"`
}

type RunMessage struct {
	CreateUserUuid string    `json:"create_user_uuid" bson:"create_user_uuid"`
	CreateAt       time.Time `json:"create_at" bson:"create_at, omitempty"`
	Message        string    `json:"message" bson:"message"`
}
type RunResponse struct {
	Time     float64 `json:"time" bson:"time, omitempty"`
	Distance float64 `json:"distance" bson:"distance, omitempty"`
}

type RunStatus int

const (
	RunActivate RunStatus = 0
	RunVerify   RunStatus = 1
	RunRetry    RunStatus = 2
)

type RunStatusResponse struct {
	Status  int    `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
