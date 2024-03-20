package models

import "time"

type Result struct {
	Questions []QuestionWithAnswersWithOutIsCorrect `json:"questions"`
}

type ResultStudentIn struct {
	Mark     int `json:"mark"`
	Score    int `json:"score"`
	MaxScore int `json:"max_score"`
}

type ResultStudent struct {
	Mark         int       `json:"mark"`
	Score        int       `json:"score"`
	MaxScore     int       `json:"max_score"`
	DateTimePass time.Time `json:"date_time_pass"`
	PassID       int       `json:"pass_id"`
	AccessID     int       `json:"access_id"`
	StudentID    int       `json:"student_id"`
	TimePass     int       `json:"time_pass"`
}
