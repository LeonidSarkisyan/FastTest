package models

type AccessOut struct {
	ID          int    `json:"id"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	TestID      int    `json:"test_id"`
	GroupID     int    `json:"group_id"`
	Shuffle     bool   `json:"shuffle"`
	PassageTime int    `json:"passage_time" binding:"required"`

	Criteria     `json:"criteria" binding:"required"`
	CriteriaJson []byte
}

type Access struct {
	Shuffle      bool   `json:"shuffle"`
	PassageTime  int    `json:"passage_time" binding:"required"`
	DateStart    string `json:"date_start" binding:"required"`
	DateEnd      string `json:"date_end" binding:"required"`
	Criteria     `json:"criteria" binding:"required"`
	CriteriaJson []byte
}

type Criteria struct {
	Five  int `json:"five" binding:"required"`
	Four  int `json:"four" binding:"required"`
	Three int `json:"three" binding:"required"`
}
