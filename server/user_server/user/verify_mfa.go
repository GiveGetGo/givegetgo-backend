package user

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"user_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// VerifyMFARequest is the request for verifying MFA
type VerifyMFARequest struct {
	UserID           uint   `json:"userID" binding:"required"`
	VerificationCode string `json:"verification_code" binding:"required"`
}

// MFAVerifyHandler handles the MFA code verification requests.
func VerifyMFAHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifyMFARequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40001",
				Message: "Invalid request body",
			})
			return
		}

		// Retrieve the user's encrypted TOTP secret from the database.
		user, err := userUtils.GetUserByID(req.UserID)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// Decode the user's encrypted TOTP secret.
		keyHex := os.Getenv("MFA_SECRET_KEY")
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		encryptedSecret, err := base64.StdEncoding.DecodeString(user.MFASecret)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// Decrypt the TOTP secret.
		decryptedSecret, err := Decrypt(encryptedSecret, key)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// Validate the TOTP code.
		isValid := totp.Validate(req.VerificationCode, string(decryptedSecret))
		if !isValid {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40005",
				Message: "Invalid MFA code",
			})
			return
		}

		// Optionally, update the user's MFA verification status in the database here.
		c.JSON(http.StatusOK, GeneralUserResponse{
			Code:    "20001",
			Message: "MFA code verified successfully",
		})
	}
}
