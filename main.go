package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/d-kuznetsov/chat/config"
	// "github.com/d-kuznetsov/chat/logger"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there !!!, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
