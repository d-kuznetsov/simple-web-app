package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/d-kuznetsov/chat/config"
	"github.com/d-kuznetsov/chat/db"
)

func main() {
	db.Connect()
	defer db.Close()

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
		renderTemplate(w, "home", nil)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username, password := r.FormValue("username"), r.FormValue("password")
			user, err := db.FindUserByName(username)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			if user != nil && user.Password == password {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			renderTemplate(w, "credentials", LoginDetails{
				"Log In",
				"/login",
				username,
				password,
				"Username or password is incorrect",
			})
			return
		}

		renderTemplate(w, "credentials", LoginDetails{Label: "Log In", Action: "/login"})
	})

	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			username, password := r.FormValue("username"), r.FormValue("password")
			user, err := db.FindUserByName(username)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			if user != nil {
				renderTemplate(w, "credentials", LoginDetails{
					"Sign Up",
					"/signup",
					username,
					password,
					"User with this username already exists",
				})
				return
			}
			db.CreateUser(username, password)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		renderTemplate(w, "credentials", LoginDetails{
			Label:  "Sign Up",
			Action: "/signup",
		})
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

type LoginDetails struct {
	Label    string
	Action   string
	Username string
	Password string
	Error    string
}
