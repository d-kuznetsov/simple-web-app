package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/d-kuznetsov/chat/db"
	"github.com/d-kuznetsov/chat/session"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !session.IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		articles, err := db.GetAllArticles()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		RenderTemplate(w, "articles", ArticleTmplOptions{
			articles,
			LayoutTmplOptions{IsAuthorized: true},
		})
	}).Methods("GET")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		username, password := r.FormValue("username"), r.FormValue("password")
		user, err := db.FindUserByName(username)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		if user != nil && user.Password == password {
			session.Login(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
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

	}).Methods("POST")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if session.IsAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		RenderTemplate(w, "credentials", CredentialsTmplOptions{
			Label:             "Log In",
			Action:            "/login",
			LayoutTmplOptions: LayoutTmplOptions{IsAuthorized: false},
		})
	}).Methods("GET")

	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		username, password := r.FormValue("username"), r.FormValue("password")
		user, err := db.FindUserByName(username)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
		db.CreateUser(username, password)
		session.Login(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods("POST")

	router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if session.IsAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		RenderTemplate(w, "credentials", CredentialsTmplOptions{
			Label:  "Sign Up",
			Action: "/signup",
			LayoutTmplOptions: LayoutTmplOptions{
				IsAuthorized: false,
			},
		})
	}).Methods("GET")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session.Logout(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	return router
}
