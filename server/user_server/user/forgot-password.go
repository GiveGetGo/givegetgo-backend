package user

import (
	"errors"
	"net/http"
	"regexp"
	"user_server/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserForgetPassRequest is the request body for the user forget password endpoint
type UserForgetPassRequest struct {
	Email string `json:"email" binding:"required"`
}

// UserForgetPassResponse is the response body for the user forget password endpoint
type UserForgetPassResponse struct {
	Event   string `json:"event"`
	Code    string `json:"code"`
	Message string `json:"message"`
}


func forgotPasswordHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserForgetPassRequest
		// Parse and validate the request body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, UserForgetPassResponse{
				Event:   "forgot-password",
				Code:    "40001",
				Message: "Invalid request format or missing data.",
			})
			return
		}
		
		// Check if the email is in the correct format
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+@purdue\.edu$`, req.Email)
		if !matched {
			c.JSON(http.StatusBadRequest, UserForgetPassResponse{
				Event:   "forgot-password",
				Code:    "40002",
				Message: "Email must be a @purdue.edu address",
			})
			return
		}

		// Check if the email exists in your database
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Optionally, you could choose not to reveal whether an email is registered
				c.JSON(http.StatusUnauthorized, UserForgetPassResponse{
					Event:   "forgot-password",
					Code:    "40101",
					Message: "Unauthorized access.",
				})
				return
			}

			// Handle internal server error
			c.JSON(http.StatusInternalServerError, UserForgetPassResponse{
				Event:   "forgot-password",
				Code:    "50001",
				Message: "Internal Server Error",
			})
			return
		}

		// call /verification/request-email endpoint in the verification server to send the verification email
		// if the verification server is down, return an internal server error
		err = userUtils.RequestForgetpassVerificationEmail(user.UserID, user.UserName, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserForgetPassResponse{
				Event:   "forgot-password",
				Code:    "50004",
				Message: "Failed to request verification email",
			})
			return
		}

		c.JSON(http.StatusOK, UserForgetPassResponse{
			Event:   "forgot-password",
			Code:    "20002",
			Message: "User forgot password",
		})
	}
}

