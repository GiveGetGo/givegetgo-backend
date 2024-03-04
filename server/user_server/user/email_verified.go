package user

import (
	"log"
	"net/http"
	"user_server/schema"
	"user_server/utils"

	"github.com/gin-gonic/gin"
)

// SetUserEmailVerifiedRequest - request body for setting user email to verified
type SetUserEmailVerifiedRequest struct {
	Email string `json:"email" binding:"required"`
}

// SetUserEmailVerifiedResponse - response body for setting user email to verified
type SetUserEmailVerifiedResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// successful email verification from verification_server, call this endpoint to set the user's email to verified
func SetUserEmailVerifiedHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req SetUserEmailVerifiedRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40001",
				Message: "Invalid request body",
			})
			return
		}

		// Update the user's email to verified
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// If the user does not exist, return an error
		if (schema.User{}) == user {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40002",
				Message: "User not found",
			})
			return
		}

		// If the user's email is already verified, return an error
		if user.EmailVerified {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40003",
				Message: "Email already verified",
			})
			return
		}

		// Update the user's email to verified
		err = userUtils.SetUserEmailVerified(req.Email)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50002",
				Message: "Failed to set email verified",
			})
			return
		}

		// Return success
		c.JSON(http.StatusOK, SetUserEmailVerifiedResponse{
			Code:    "20005",
			Message: "Email verified",
		})
	}
}
