package models

import "time"

type PassesIn struct {
	Code      int64
	StudentID int
}

type Passes struct {
	ID                int
	DateTimeActivated *time.Time `json:"date_time_activated"`
	IsActivated       bool       `json:"is_activated"`
	Code              int64      `json:"code"`
	StudentID         int        `json:"student_id"`
}

type PassesOut struct {
	ID        int
	StudentID int `json:"student_id"`
}
