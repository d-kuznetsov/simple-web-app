package router

import (
	"net/http"

	"github.com/d-kuznetsov/chat/adapter"
	"github.com/d-kuznetsov/chat/session"
)

func LogInGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if isAuthenticated {
		RedirectToArticles(w, r)
		return
	}
	data := struct{ LayoutTmplOpts }{LayoutTmplOpts: LayoutTmplOpts{IsAuthorized: false}}
	RenderTemplate(w, "login", data)
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
	data := struct{ LayoutTmplOpts }{
		LayoutTmplOpts: LayoutTmplOpts{IsAuthorized: false, Error: "Username or password is incorrect"},
	}
	RenderTemplate(w, "login", data)
}

func SignUpGetHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated, _ := session.IsAuthenticated(r)
	if isAuthenticated {
		RedirectToArticles(w, r)
		return
	}
	data := struct{ LayoutTmplOpts }{
		LayoutTmplOpts: LayoutTmplOpts{IsAuthorized: false},
	}
	RenderTemplate(w, "signup", data)
}

func SignUpPostHandler(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	user, err := adapter.FindUserByName(username)
	if err != nil {
		ThrowServerError(w)
		return
	}
	if user != nil {
		data := struct{ LayoutTmplOpts }{
			LayoutTmplOpts: LayoutTmplOpts{IsAuthorized: false, Error: "User with this username already exists"},
		}
		RenderTemplate(w, "signup", data)
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
