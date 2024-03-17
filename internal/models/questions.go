package models

type Question struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type QuestionWithAnswers struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Answers []Answer `json:"answers"`
}

type QuestionUpdate struct {
	Text string `json:"text"`
}
