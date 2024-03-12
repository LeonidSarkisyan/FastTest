package models

type UserIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserOut struct {
	ID    int
	Email string
}

type User struct {
	ID       int
	Email    string
	Password string
}
