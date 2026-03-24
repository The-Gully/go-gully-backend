package auth

import (
	"encoding/gob"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const SessionKey = "user_id"
const SessionName = "session"

var store *sessions.CookieStore

func init() {
	gob.Register(uint(0))
}

func InitSession(sessionSecret string) {
	store = sessions.NewCookieStore([]byte(sessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   os.Getenv("GIN_MODE") == "release",
		SameSite: http.SameSiteLaxMode,
	}
}

func GetStore() *sessions.CookieStore {
	return store
}
