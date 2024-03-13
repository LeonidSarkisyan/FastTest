package models

type TestIn struct {
	Title string `json:"title"`
}

type Test struct {
	Title  string
	UserID int
}

type TestOut struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	DateTimeCreate int64  `json:"date_time_create"`
}

type TestUpdate struct {
	Title string `json:"title"`
}
