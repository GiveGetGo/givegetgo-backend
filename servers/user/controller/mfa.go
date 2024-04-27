package controller

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"user/utils"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
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

		// get user from session
		session := sessions.Default(c)
		userIdValue := session.Get("userid")
		if userIdValue == nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		userId, ok := userIdValue.(uint)
		if !ok {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		user, err := userUtils.GetUserByID(uint(userId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.ResponseError(c, http.StatusNotFound, types.UserNotFound())
			} else {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			}
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
		err = userUtils.StoreEncryptedTOTPSecret(user.UserID, encryptedSecretBase64)
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

		// get user from session
		session := sessions.Default(c)
		userIdValue := session.Get("userid")
		if userIdValue == nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		userId, ok := userIdValue.(uint)
		if !ok {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		user, err := userUtils.GetUserByID(uint(userId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.ResponseError(c, http.StatusNotFound, types.UserNotFound())
			} else {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			}
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
		err = userUtils.MarkMFAVerified(user.UserID)
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
		session := sessions.Default(c)
		userIdValue := session.Get("userid")
		if userIdValue == nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		userId, ok := userIdValue.(uint)
		if !ok {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		user, err := userUtils.GetUserByID(uint(userId))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.ResponseError(c, http.StatusNotFound, types.UserNotFound())
			} else {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			}
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
