package router

import (
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", ArticlesGetHandler).Methods("GET")
	router.HandleFunc("/login", LogInPostHandler).Methods("POST")
	router.HandleFunc("/login", LogInGetHandler).Methods("GET")
	router.HandleFunc("/signup", SignUpPostHandler).Methods("POST")
	router.HandleFunc("/signup", SignUpGetHandler).Methods("GET")
	router.HandleFunc("/logout", LogOutHandler)
	router.HandleFunc("/articles/{id}", OneArticleGetHandler).Methods("GET")
	router.HandleFunc("/create-article", CreateArticleGetHandler).Methods("GET")
	router.HandleFunc("/create-article", CreateArticlePostHandler).Methods("POST")
	router.HandleFunc("/update-article/{id}", UpdateArticleGetHandler).Methods("GET")
	router.HandleFunc("/update-article", UpdateArticlePostHandler).Methods("POST")
	router.HandleFunc("/my-articles", ArticlesOfUserGetHandler).Methods("GET")
	router.HandleFunc("/delete-articles", DeleteArticlesPostHandler).Methods("POST")

	return router
}
