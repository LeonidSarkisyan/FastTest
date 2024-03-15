package models

type GroupIn struct {
	Title string `json:"title"`
}

type GroupUpdate struct {
	Title string `json:"title"`
}

type GroupOut struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}
