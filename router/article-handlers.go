package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/d-kuznetsov/chat/adapter"
	"github.com/d-kuznetsov/chat/models"
	"github.com/d-kuznetsov/chat/session"
)

func ArticlesGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	articles, err := adapter.GetAllArticles()
	if err != nil {
		ThrowServerError(w)
		return
	}
	data := struct {
		Articles []models.Article
		LayoutTmplOpts
	}{
		articles,
		LayoutTmplOpts{IsAuthorized: true},
	}
	RenderTemplate(w, "articles", data)
}

func OneArticleGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	vars := mux.Vars(r)
	article, err := adapter.GetArticleById(vars["id"])
	if err != nil {
		ThrowServerError(w)
		return
	}
	data := struct {
		models.Article
		LayoutTmplOpts
	}{
		*article,
		LayoutTmplOpts{IsAuthorized: true},
	}
	RenderTemplate(w, "article", data)
}

func CreateArticleGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	data := struct{ LayoutTmplOpts }{LayoutTmplOpts{IsAuthorized: true}}
	RenderTemplate(w, "createArticle", data)
}

func CreateArticlePostHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, userId := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	title, text := r.FormValue("title"), r.FormValue("text")
	_, err := adapter.CreateArticle(title, text, userId)
	if err != nil {
		ThrowServerError(w)
		return
	}
	http.Redirect(w, r, "/create-article", http.StatusFound)
}

func UpdateArticleGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	vars := mux.Vars(r)
	article, err := adapter.GetArticleById(vars["id"])
	if err != nil {
		ThrowServerError(w)
		return
	}
	data := struct {
		models.Article
		LayoutTmplOpts
	}{
		*article,
		LayoutTmplOpts{IsAuthorized: true},
	}
	RenderTemplate(w, "updateArticle", data)
}

func UpdateArticlePostHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, userId := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	id, title, text := r.FormValue("id"), r.FormValue("title"), r.FormValue("text")
	article, _ := adapter.GetArticleById(id)
	if article.User != userId {
		ThrowServerError(w)
		return
	}
	err := adapter.UpdateArticle(id, title, text)
	if err != nil {
		ThrowServerError(w)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	RedirectToArticles(w, r)
}

func ArticlesOfUserGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, userId := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	articles, err := adapter.GetArticlesByUserId(userId)
	if err != nil {
		ThrowServerError(w)
		return
	}
	data := struct {
		Articles []models.Article
		LayoutTmplOpts
	}{
		articles,
		LayoutTmplOpts{IsAuthorized: true},
	}
	RenderTemplate(w, "myArticles", data)
}

func DeleteArticlesPostHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	r.ParseForm()
	err := adapter.DeleteArticlesByIds(r.Form["articles"])
	if err != nil {
		ThrowServerError(w)
		return
	}
	http.Redirect(w, r, "/my-articles", http.StatusFound)
}
