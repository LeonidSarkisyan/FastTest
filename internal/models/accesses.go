package models

type Access struct {
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
