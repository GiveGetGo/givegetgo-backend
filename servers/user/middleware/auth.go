package middleware

import (
	"log"
	"net/http"
	time "time"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("userid") == nil {
			log.Println("Unauthorized")
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			c.Abort()
			return
		}
		c.Next()
		// We will reset the expiration time on every request after the user is authenticated
		cookie, _ := c.Request.Cookie("givegetgo")
		if cookie != nil {
			cookie.Expires = time.Now().UTC().Add(time.Hour * 24 * 30)
			http.SetCookie(c.Writer, cookie)
		}
	}
}
