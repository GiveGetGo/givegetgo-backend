package user

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"user_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// func MFAQRCodeHandler - generate a QR code for MFA
func MFAQRCodeHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userid") // Assuming you're getting the userID from the URL
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, GeneralUserResponse{
				Code:    "40001",
				Message: "Invalid Request Body",
			})
			return
		}

		// Get the user
		user, err := userUtils.GetUserByID(uint(userIDUint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// Retrieve the Base64-encoded encrypted TOTP secret and decode it
		encryptedSecretBase64 := user.MFASecret
		encryptedSecret, err := base64.StdEncoding.DecodeString(encryptedSecretBase64)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// Decrypt the secret
		keyHex := os.Getenv("MFA_SECRET_KEY")
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		decryptedSecret, err := Decrypt(encryptedSecret, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// Generate the URI for the QR code
		uri := fmt.Sprintf("otpauth://totp/GiveGetGo:%s?secret=%s&issuer=GiveGetGo", user.Email, string(decryptedSecret))

		// Generate the QR code
		qrCode, err := qrcode.Encode(uri, qrcode.Medium, 256)
		if err != nil {
			c.JSON(http.StatusInternalServerError, GeneralUserResponse{
				Code:    "50001",
				Message: "Internal server error",
			})
			return
		}

		// Respond with the QR code image
		c.Data(http.StatusOK, "image/png", qrCode)
	}
}
