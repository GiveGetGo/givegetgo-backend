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

// MFAVerificationRequest
type MFAVerificationRequest struct {
	UserID uint   `json:"userID" binding:"required"`
	Email  string `json:"email" binding:"required"`
}

// func RequestMFAVerificationHandler - request MFA verification
func RequestMFAVerificationHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req MFAVerificationRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40001",
				Message: "Invalid request body",
			})
			return
		}

		// Check if user exists
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40002",
				Message: "User not found",
			})
			return
		}

		// Check if the user's email is verified
		if !user.EmailVerified {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40003",
				Message: "Email not verified",
			})
			return
		}

		// Check if the user is already MFA verified
		if user.MFAVerified {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40004",
				Message: "MFA already verified",
			})
			return
		}

		// generate a TOTP secret
		secret, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "GiveGetGo",
			AccountName: req.Email,
		})
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// encrypt the secret
		key, err := hex.DecodeString(os.Getenv("MFA_SECRET_KEY"))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		encryptedSecret, err := Encrypt([]byte(secret.Secret()), key)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// store the encrypted secret in the database
		encryptedSecretBase64 := base64.StdEncoding.EncodeToString(encryptedSecret)
		err = userUtils.StoreEncryptedTOTPSecret(req.UserID, encryptedSecretBase64)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// On success, return the QR code
		c.JSON(http.StatusOK, GeneralUserResponse{
			Code:    "20006",
			Message: "MFA verification requested",
		})
	}
}
