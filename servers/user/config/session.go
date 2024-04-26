package config

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

// InitSession initializes the session store with configuration from Viper
func InitSession() sessions.Store {
	sessionKey := config.GetString("SESSION_SECRET")
	log.Printf("session: %s", sessionKey)

	store := cookie.NewStore([]byte(sessionKey))
	store.Options(sessions.Options{
		Path:     config.GetString("session.path"),
		MaxAge:   config.GetInt("session.max_age"),
		Secure:   config.GetBool("session.secure"),
		HttpOnly: config.GetBool("session.http_only"),
		Domain:   config.GetString("session.domain"),
	})

	return store
}
