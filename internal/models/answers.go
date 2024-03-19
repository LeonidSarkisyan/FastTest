package models

type Answer struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"-"`
}

type AnswerUpdate struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}
