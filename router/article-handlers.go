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
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	vars := mux.Vars(r)
	article, err := adapter.GetArticleById(vars["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	RenderTemplate(w, "createArticle", LayoutTmplOptions{IsAuthorized: true})
}

func CreateArticlePostHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, userId := session.IsAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	title, text := r.FormValue("title"), r.FormValue("text")
	_, err := adapter.CreateArticle(title, text, userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func UpdateArticleGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	vars := mux.Vars(r)
	a, err := adapter.GetArticleById(vars["id"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	RenderTemplate(w, "updateArticle", OneArticleTmplOptions{
		*a,
		LayoutTmplOptions{IsAuthorized: true},
	})
}

func UpdateArticlePostHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, userId := session.IsAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	id, title, text := r.FormValue("id"), r.FormValue("title"), r.FormValue("text")
	a, _ := adapter.GetArticleById(id)
	if a.User != userId {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	err := adapter.UpdateArticle(id, title, text)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func ArticlesOfUserGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, userId := session.IsAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	articles, err := adapter.GetArticlesByUserId(userId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	RenderTemplate(w, "articlesOfUser", ArticleTmplOptions{
		articles,
		LayoutTmplOptions{IsAuthorized: true},
	})
}

func DeleteArticlesPostHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if !isAuthenticated {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	r.ParseForm()
	err := adapter.DeleteArticlesByIds(r.Form["articles"])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/my-articles", http.StatusFound)
}
