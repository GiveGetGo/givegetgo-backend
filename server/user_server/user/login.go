package user

import (
	"errors"
	"net/http"
	"user_server/schema"
	"user_server/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Event   string `json:"event"`
	Code    string `json:"code"`
	UserID  uint   `json:"user_id"`
	Message string `json:"msg"`
}

func authenticate(user schema.User, password string) bool {
	// Compare the hash of the input password with the hash stored in the database.
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	return err == nil
}

func LoginHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40001",
				Message: "Invalid request format or missing data.",
			})
			return
		}

		// Check if the email is already registered
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// If email not verified, return an error
		if !user.EmailVerified {
			c.JSON(http.StatusUnauthorized, GeneralUserResponse{
				Code:    "40102",
				Message: "Email not verified",
			})
			return
		}

		authenticated := authenticate(user, req.Password)
		if authenticated {
			c.JSON(http.StatusOK, LoginResponse{
				Event:   "login",
				Code:    "20001",
				Message: "Wait for MFA verification",
			})
		} else {
			c.JSON(http.StatusUnauthorized, GeneralUserResponse{
				Code:    "40101",
				Message: "Invalid credentials",
			})
		}
	}
}
