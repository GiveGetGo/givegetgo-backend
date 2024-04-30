package controller

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"user/utils"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

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

		res.ResponseSuccess(c, http.StatusOK, "verify-mfa", types.Success())
	}
}

func GetMFAHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
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

		keyHex := os.Getenv("MFA_SECRET_KEY")
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Check if MFA is already set up
		if user.MFASecret != "" {
			// Decrypt and generate QR code from existing secret
			encryptedSecret, err := base64.StdEncoding.DecodeString(user.MFASecret)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			decryptedSecret, err := utils.Decrypt(encryptedSecret, key)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			userUtils.GenerateAndSendQRCode(c, user.Email, decryptedSecret)
		} else {
			// Generate a new secret and store it encrypted
			secret, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "GiveGetGo",
				AccountName: user.Email,
			})
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			encryptedSecret, err := utils.Encrypt([]byte(secret.Secret()), key)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			encryptedSecretBase64 := base64.StdEncoding.EncodeToString(encryptedSecret)
			err = userUtils.StoreEncryptedTOTPSecret(user.UserID, encryptedSecretBase64)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			userUtils.GenerateAndSendQRCode(c, user.Email, []byte(secret.Secret()))
		}
	}
}
