package session

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/d-kuznetsov/chat/config"
)

var store *sessions.CookieStore
var sessionName = "session-name"

func CreateStore() {
	store = sessions.NewCookieStore([]byte(config.SessionKey))
}

func IsAuthenticated(r *http.Request) bool {
	checkStore()
	session, _ := store.Get(r, sessionName)
	auth, ok := session.Values["authenticated"].(bool)
	fmt.Println(auth, ok)
	return auth && ok
}

func Login(w http.ResponseWriter, r *http.Request) {
	checkStore()
	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	checkStore()
	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func checkStore() {
	if store == nil {
		log.Fatal("Store does not exist")
	}
}
