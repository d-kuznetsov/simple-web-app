package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/d-kuznetsov/chat/config"
)

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles("templates/home.html", "templates/layout.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
