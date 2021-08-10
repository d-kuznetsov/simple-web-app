package router

import (
	"html/template"
	"log"
	"net/http"

	"github.com/d-kuznetsov/chat/models"
)

var Templates = map[string]*template.Template{
	"credentials": template.Must(template.ParseFiles("templates/layout.html", "templates/credentials.html")),
	"home":        template.Must(template.ParseFiles("templates/layout.html", "templates/home.html")),
}

var RenderTemplate = func(w http.ResponseWriter, tmplName string, tmplData interface{}) {
	tmpl, ok := Templates[tmplName]
	if !ok {
		log.Fatal("There is not the template")
	}
	err := tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

type ArticleTmplOptions struct {
	Articles []models.Article
	LayoutTmplOptions
}
