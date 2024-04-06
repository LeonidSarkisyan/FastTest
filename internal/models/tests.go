package models

type TestIn struct {
	Title string `json:"title"`
}

type Test struct {
	Title  string `json:"title"`
	UserID int
}

type TestOut struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	DateTimeCreate int64  `json:"date_time_create"`
	Count          int    `json:"count"`
	IsDeleted      bool   `json:"is_deleted"`
}

type TestUpdate struct {
	Title string `json:"title"`
}
