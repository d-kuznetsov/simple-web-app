package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/d-kuznetsov/chat/config"
	"github.com/d-kuznetsov/chat/db"
	"github.com/d-kuznetsov/chat/session"
)

func main() {
	db.Connect()
	defer db.Close()
	session.CreateStore()

	templates := map[string]*template.Template{
		"credentials": template.Must(template.ParseFiles("templates/layout.html", "templates/credentials.html")),
		"home":        template.Must(template.ParseFiles("templates/layout.html", "templates/home.html")),
	}
	renderTemplate := func(w http.ResponseWriter, tmplName string, tmplData interface{}) {
		tmpl, ok := templates[tmplName]
		if !ok {
			log.Fatal("There is not the template")
		}
		err := tmpl.Execute(w, tmplData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !session.IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		renderTemplate(w, "home", LayoutTmplOptions{IsAuthorized: true})
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
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
			renderTemplate(w, "credentials", CredentialsTmplOptions{
				"Log In",
				"/login",
				username,
				password,
				"Username or password is incorrect",
				LayoutTmplOptions{IsAuthorized: false},
			})
			return
		}
		if session.IsAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		renderTemplate(w, "credentials", CredentialsTmplOptions{
			Label:             "Log In",
			Action:            "/login",
			LayoutTmplOptions: LayoutTmplOptions{IsAuthorized: false},
		})
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username, password := r.FormValue("username"), r.FormValue("password")
			user, err := db.FindUserByName(username)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			if user != nil {
				renderTemplate(w, "credentials", CredentialsTmplOptions{
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
			return
		}
		if session.IsAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		renderTemplate(w, "credentials", CredentialsTmplOptions{
			Label:  "Sign Up",
			Action: "/signup",
			LayoutTmplOptions: LayoutTmplOptions{
				IsAuthorized: false,
			},
		})
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		session.Logout(w, r)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

type LayoutTmplOptions struct {
	IsAuthorized bool
}
type CredentialsTmplOptions struct {
	Label    string
	Action   string
	Username string
	Password string
	Error    string
	LayoutTmplOptions
}
