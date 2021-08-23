package session

import (
	"log"
	"net/http"

	"github.com/d-kuznetsov/blog/config"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore
var sessionName = "session-name"

func CreateStore() {
	store = sessions.NewCookieStore([]byte(config.SessionKey))
}

func IsAuthenticated(r *http.Request) (bool, string) {
	checkStore()
	session, _ := store.Get(r, sessionName)
	auth, ok := session.Values["authenticated"].(bool)
	var id string
	if session.Values["user-id"] == nil {
		id = ""
	} else {
		id = session.Values["user-id"].(string)
	}
	return auth && ok, id
}

func Login(w http.ResponseWriter, r *http.Request, id string) {
	checkStore()
	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = true
	session.Values["user-id"] = id
	session.Save(r, w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	checkStore()
	session, _ := store.Get(r, sessionName)
	session.Values["authenticated"] = false
	session.Values["user-id"] = ""
	session.Save(r, w)
}

func checkStore() {
	if store == nil {
		log.Fatal("Store does not exist")
	}
}
