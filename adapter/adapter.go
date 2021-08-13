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

func CreateUser(username, password string) (string, error) {
	objId, err := db.CreateUser(username, password)
	return objId.Hex(), err
}

func GetAllArticles() ([]models.Article, error) {
	arts, err := db.GetAllArticles()
	var articles []models.Article
	for _, val := range arts {
		articles = append(articles, *convertArticle(&val))
	}
	return articles, err
}

func GetArticleById(id string) (*models.Article, error) {
	a, err := db.GetArticleById(id)
	if err != nil {
		return nil, err
	}
	return convertArticle(a), err
}

func convertArticle(a *db.Article) *models.Article {
	if a == nil {
		return nil
	}
	return &models.Article{
		Id:    a.Id.Hex(),
		Title: a.Title,
		Date:  a.Date,
		Text:  a.Text,
		User:  a.User.Hex(),
	}
}

func CreateArticle(title, text, userId string) error {
	_, err := db.CreateArticle(title, text, userId)
	return err
}
