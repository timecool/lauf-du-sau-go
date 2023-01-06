package models

import (
	"time"
)

type Run struct {
	UUID     string    `json:"uuid" bson:"_id, omitempty"`
	Date     time.Time `json:"date" bson:"date, omitempty"`
	CreateAt time.Time `json:"createAt" bson:"createAt, omitempty"`
	Time     float64   `json:"time" bson:"time, omitempty"`
	Distance float64   `json:"distance" bson:"distance, omitempty"`
	Url      string    `json:"url" bson:"url, omitempty"`
	Status   RunStatus `json:"status" bson:"status, omitempty"`
}
type RunResponse struct {
	Time     float64 `json:"time" bson:"time, omitempty"`
	Distance float64 `json:"distance" bson:"distance, omitempty"`
}

type RunStatus int

const (
	Activate RunStatus = 0
	Verify   RunStatus = 1
	Retry    RunStatus = 2
)
