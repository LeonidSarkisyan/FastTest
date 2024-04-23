package models

type Question struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Data any    `json:"data"`
	Type string `json:"type"`
}

type QuestionWithAnswers struct {
	ID       int      `json:"id"`
	Text     string   `json:"text"`
	Type     string   `json:"type"`
	Data     any      `json:"data"`
	ImageURL string   `json:"image_url"`
	Answers  []Answer `json:"answers"`
}

type QuestionsFromChatGPT struct {
	Questions []QuestionWithAnswersWithOutIsCorrect `json:"questions"`
}

type QuestionWithAnswersWithOutIsCorrect struct {
	ID       int                 `json:"id"`
	Text     string              `json:"text"`
	Type     string              `json:"type"`
	Data     any                 `json:"data"`
	ImageURL string              `json:"image_url"`
	Answers  []AnswerWithCorrect `json:"answers"`
}

type QuestionUpdate struct {
	Text string `json:"text"`
}

type QuestionForMap struct {
	Question
	AnswersMap map[int]Answer
}

type QuestionGroupData struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Name    string   `json:"name"`
	Answers []string `json:"answers"`
}

type QuestionRangeData struct {
	Ranges []Range `json:"ranges"`
}

type Range struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}
