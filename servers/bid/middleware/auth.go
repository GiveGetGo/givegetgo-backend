package middleware

import (
	"net/http"
	"os"
	time "time"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the session cookie from the request
		cookie, err := c.Request.Cookie("givegetgo")
		if err != nil || cookie.Value == "" {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			c.Abort()
			return
		}

		// Prepare the request to verify the session via the user service
		userServiceURL := os.Getenv("USER_SERVICE_URL") + "/v1/user/session"
		client := &http.Client{Timeout: 10 * time.Second}
		req, _ := http.NewRequest("GET", userServiceURL, nil)

		// Make the HTTP request to the user service
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			c.Abort()
			return
		}

		// Assuming the session is valid, proceed with the request and refresh the session
		c.Next()

		// reset the expiration time on every request after the user is authenticated
		cookie.Expires = time.Now().UTC().Add(time.Hour * 24) // Set cookie to expire in 1 day
		http.SetCookie(c.Writer, cookie)
	}
}
