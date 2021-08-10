package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d-kuznetsov/chat/config"
	"github.com/d-kuznetsov/chat/db"
	"github.com/d-kuznetsov/chat/router"
	"github.com/d-kuznetsov/chat/session"
)

func main() {
	db.Connect()
	defer db.Close()

	session.CreateStore()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":"+config.Port, router.GetRouter()))
}
