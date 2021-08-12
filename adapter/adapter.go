package adapter

import (
	"github.com/d-kuznetsov/chat/db"
	"github.com/d-kuznetsov/chat/models"
)

func FindUserByName(name string) (*models.User, error) {
	u, err := db.FindUserByName(name)
	return &models.User{
		Id:       u.Id.Hex(),
		Username: u.Username,
		Password: u.Password,
	}, err
}

func CreateUser(username, password string) error {
	_, err := db.CreateUser(username, password)
	return err
}

func GetAllArticles() ([]models.Article, error) {
	arts, err := db.GetAllArticles()
	var articles []models.Article
	for _, val := range arts {
		articles = append(articles, models.Article{
			Id:    val.Id.Hex(),
			Title: val.Title,
			Date:  val.Date,
			Text:  val.Text,
			User:  val.User.Hex(),
		})
	}
	return articles, err
}
