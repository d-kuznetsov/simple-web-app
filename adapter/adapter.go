package adapter

import (
	"github.com/d-kuznetsov/blog/db"
	"github.com/d-kuznetsov/blog/models"
)

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

func FindUserByName(name string) (*models.User, error) {
	dbUser, err := db.FindUserByName(name)
	if err != nil || dbUser == nil {
		return nil, err
	}
	return &models.User{
		Id:       dbUser.Id.Hex(),
		Username: dbUser.Username,
		Password: dbUser.Password,
	}, err
}

func CreateUser(username, password string) (string, error) {
	objId, err := db.CreateUser(username, password)
	return objId.Hex(), err
}

func GetAllArticles() ([]*models.Article, error) {
	dbArticles, err := db.GetAllArticles()
	var articles []*models.Article
	for _, val := range dbArticles {
		articles = append(articles, convertArticle(val))
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

func CreateArticle(title, text, userId string) (string, error) {
	objId, err := db.CreateArticle(title, text, userId)
	return objId.Hex(), err
}

func UpdateArticle(id, title, text string) error {
	_, err := db.UpdateArticle(id, title, text)
	return err
}

func GetArticlesByUserId(userId string) ([]*models.Article, error) {
	dbArticles, err := db.GetArticlesByUserId(userId)
	var articles []*models.Article
	for _, val := range dbArticles {
		articles = append(articles, convertArticle(val))
	}
	return articles, err
}

func DeleteArticlesByIds(ids []string) error {
	_, err := db.DeleteArticlesByIds(ids)
	return err
}
