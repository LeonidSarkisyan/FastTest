package models

import "time"

type ForResultTable struct {
	Students []Student       `json:"students"`
	Passes   []Passes        `json:"passes"`
	Results  []ResultStudent `json:"results"`
}

type PassesIn struct {
	Code      int64
	StudentID int
}

type Passes struct {
	ID                int        `json:"id"`
	DateTimeActivated *time.Time `json:"date_time_activated"`
	IsActivated       bool       `json:"is_activated"`
	Code              int64      `json:"code"`
	StudentID         int        `json:"student_id"`
}

type PassesOut struct {
	ID        int
	StudentID int `json:"student_id"`
}
