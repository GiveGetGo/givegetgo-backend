package controller

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"user/utils"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// func RequestMFAVerificationHandler - request MFA verification
func RequestMFAVerificationHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req types.MFAVerificationRequest
		if err := c.BindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Check if user exists
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {

			return
		}

		// Check if the user's email is verified
		if !user.EmailVerified {
			res.ResponseError(c, http.StatusBadRequest, types.EmailNotVerified())
			return
		}

		// Check if the user is already MFA verified
		if user.MFAVerified {
			res.ResponseError(c, http.StatusBadRequest, types.AlreadyVerified())
			return
		}

		// generate a TOTP secret
		secret, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "GiveGetGo",
			AccountName: req.Email,
		})
		if err != nil {
			log.Println(err)
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// encrypt the secret
		key, err := hex.DecodeString(os.Getenv("MFA_SECRET_KEY"))
		if err != nil {
			log.Println(err)
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		encryptedSecret, err := utils.Encrypt([]byte(secret.Secret()), key)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// store the encrypted secret in the database
		encryptedSecretBase64 := base64.StdEncoding.EncodeToString(encryptedSecret)
		err = userUtils.StoreEncryptedTOTPSecret(req.UserID, encryptedSecretBase64)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "request-mfa", types.Success())
	}
}

// MFAVerifyHandler handles the MFA code verification requests.
func VerifyMFAHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.VerifyMFARequest
		if err := c.BindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Retrieve the user's encrypted TOTP secret from the database.
		user, err := userUtils.GetUserByID(req.UserID)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Decode the user's encrypted TOTP secret.
		keyHex := os.Getenv("MFA_SECRET_KEY")
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		encryptedSecret, err := base64.StdEncoding.DecodeString(user.MFASecret)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Decrypt the TOTP secret.
		decryptedSecret, err := utils.Decrypt(encryptedSecret, key)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Validate the TOTP code.
		isValid := totp.Validate(req.VerificationCode, string(decryptedSecret))
		if !isValid {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidVerification())
			return
		}

		// Mark the user as MFA verified.
		err = userUtils.MarkMFAVerified(req.UserID)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "verify-mfa", types.MFAVerified())
	}
}

// func MFAQRCodeHandler - generate a QR code for MFA
func MFAQRCodeHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userid") // Assuming you're getting the userID from the URL
		userIDUint, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Get the user
		user, err := userUtils.GetUserByID(uint(userIDUint))
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Retrieve the Base64-encoded encrypted TOTP secret and decode it
		encryptedSecretBase64 := user.MFASecret
		encryptedSecret, err := base64.StdEncoding.DecodeString(encryptedSecretBase64)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Decrypt the secret
		keyHex := os.Getenv("MFA_SECRET_KEY")
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		decryptedSecret, err := utils.Decrypt(encryptedSecret, key)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Generate the URI for the QR code
		uri := fmt.Sprintf("otpauth://totp/GiveGetGo:%s?secret=%s&issuer=GiveGetGo", user.Email, string(decryptedSecret))

		// Generate the QR code
		qrCode, err := qrcode.Encode(uri, qrcode.Medium, 256)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Respond with the QR code image
		c.Data(http.StatusOK, "image/png", qrCode)
	}
}
