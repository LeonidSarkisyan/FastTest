package models

type Question struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type QuestionUpdate struct {
	Text string `json:"text" binding:"required"`
}
