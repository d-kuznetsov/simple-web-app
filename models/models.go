package models

type User struct {
	Id       string
	Username string
	Password string
}

type Article struct {
	Id    string
	Title string
	Date  string
	Text  string
	User  string
}
