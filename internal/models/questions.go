package models

type Question struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type QuestionWithAnswers struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Type    string   `json:"type"`
	Answers []Answer `json:"answers"`
}

type QuestionWithAnswersWithOutIsCorrect struct {
	ID      int                 `json:"id"`
	Text    string              `json:"text"`
	Answers []AnswerWithCorrect `json:"answers"`
}

type QuestionUpdate struct {
	Text string `json:"text"`
}

type QuestionForMap struct {
	Question
	AnswersMap map[int]Answer
}
