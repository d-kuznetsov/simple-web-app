package router

import (
	"net/http"

	"github.com/d-kuznetsov/chat/adapter"
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
	RenderTemplate(w, "articles", ArticleTmplOptions{
		articles,
		LayoutTmplOptions{IsAuthorized: true},
	})
}

func LogInGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if isAuthenticated {
		RedirectToArticles(w, r)
		return
	}
	RenderTemplate(w, "credentials", CredentialsTmplOptions{
		Label:             "Log In",
		Action:            "/login",
		LayoutTmplOptions: LayoutTmplOptions{IsAuthorized: false},
	})
}

func LogInPostHandler(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	user, err := adapter.FindUserByName(username)
	if err != nil {
		ThrowServerError(w)
		return
	}
	if user != nil && user.Password == password {
		session.Login(w, r, user.Id)
		RedirectToArticles(w, r)
		return
	}
	RenderTemplate(w, "credentials", CredentialsTmplOptions{
		"Log In",
		"/login",
		username,
		password,
		"Username or password is incorrect",
		LayoutTmplOptions{IsAuthorized: false},
	})
}

func SignUpGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if isAuthenticated {
		RedirectToArticles(w, r)
		return
	}
	RenderTemplate(w, "credentials", CredentialsTmplOptions{
		Label:  "Sign Up",
		Action: "/signup",
		LayoutTmplOptions: LayoutTmplOptions{
			IsAuthorized: false,
		},
	})
}

func SignUpPostHandler(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	user, err := adapter.FindUserByName(username)
	if err != nil {
		ThrowServerError(w)
		return
	}
	if user != nil {
		RenderTemplate(w, "credentials", CredentialsTmplOptions{
			"Sign Up",
			"/signup",
			username,
			password,
			"User with this username already exists",
			LayoutTmplOptions{IsAuthorized: false},
		})
		return
	}
	userId, err := adapter.CreateUser(username, password)
	if err != nil {
		ThrowServerError(w)
		return
	}
	session.Login(w, r, userId)
	RedirectToArticles(w, r)
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	session.Logout(w, r)
	RedirectToArticles(w, r)
}
