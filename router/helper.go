package router

import (
	"html/template"
	"log"
	"net/http"
)

func ThrowServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func RedirectToArticles(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func RedirectToLogIn(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

var templates = map[string]*template.Template{
	"login":         template.Must(template.ParseFiles("templates/layout.html", "templates/login.html")),
	"signup":        template.Must(template.ParseFiles("templates/layout.html", "templates/signup.html")),
	"articles":      template.Must(template.ParseFiles("templates/layout.html", "templates/articles.html")),
	"article":       template.Must(template.ParseFiles("templates/layout.html", "templates/article.html")),
	"createArticle": template.Must(template.ParseFiles("templates/layout.html", "templates/create-article.html")),
	"updateArticle": template.Must(template.ParseFiles("templates/layout.html", "templates/update-article.html")),
	"myArticles":    template.Must(template.ParseFiles("templates/layout.html", "templates/my-articles.html")),
}

func RenderTemplate(w http.ResponseWriter, tmplName string, tmplData interface{}) {
	tmpl, ok := templates[tmplName]
	if !ok {
		log.Fatal("Templete does not exist")
	}
	err := tmpl.Execute(w, tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type LayoutTmplOpts struct {
	IsAuthorized bool
	Error        string
}
