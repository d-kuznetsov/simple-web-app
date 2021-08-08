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
		"login": template.Must(template.ParseFiles("templates/layout.html", "templates/login.html")),
		"home":  template.Must(template.ParseFiles("templates/layout.html", "templates/home.html")),
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
			if r.FormValue("username") == "user" && r.FormValue("password") == "1234" {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			renderTemplate(w, "login", LoginDetails{
				r.FormValue("username"),
				r.FormValue("password"),
				"Data is incorrect",
			})
			return
		}

		renderTemplate(w, "login", nil)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}

type LoginDetails struct {
	Username string
	Password string
	Error    string
}
