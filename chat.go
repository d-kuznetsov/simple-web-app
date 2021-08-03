package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there !!!, I love %s!", r.URL.Path[1:])
}

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "port")
	flag.Parse()
	fmt.Printf("%T %s", port, port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
