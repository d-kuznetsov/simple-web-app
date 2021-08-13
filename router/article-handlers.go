package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/d-kuznetsov/chat/adapter"
	"github.com/d-kuznetsov/chat/session"
)

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
	RenderTemplate(w, "article", OneArticleTmplOptions{
		*article,
		LayoutTmplOptions{IsAuthorized: true},
	})
}

func CreateArticleGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if !isAuthenticated {
		RedirectToLogIn(w, r)
		return
	}
	RenderTemplate(w, "createArticle", LayoutTmplOptions{IsAuthorized: true})
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
	a, err := adapter.GetArticleById(vars["id"])
	if err != nil {
		ThrowServerError(w)
		return
	}
	RenderTemplate(w, "updateArticle", OneArticleTmplOptions{
		*a,
		LayoutTmplOptions{IsAuthorized: true},
	})
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
	RenderTemplate(w, "articlesOfUser", ArticleTmplOptions{
		articles,
		LayoutTmplOptions{IsAuthorized: true},
	})
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
