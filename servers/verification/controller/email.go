package controller

import (
	"log"
	"net/http"
	"regexp"
	"verification/utils"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

// func RequestEmailVerificationHandler - verifiy the verification code
func RequestEmailVerificationHandler(verificationUtils utils.IVerificationUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req types.GetEmailVerificationRequest
		if err := c.BindJSON(&req); err != nil {
			log.Println("Error parsing request body: ", err)
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// get the context
		ctx := c.Request.Context()

		// Check if the email is in the correct format
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+@purdue\.edu$`, req.Email)
		if !matched {
			log.Println("Invalid email: ", req.Email)
			res.ResponseError(c, http.StatusBadRequest, types.InvalidEmail())
			return
		}

		// idnetify the event
		switch req.Event {
		case types.RegisterEvent:
			// generate a verification code
			verificationCode, err := verificationUtils.GenerateRegisterVerificationCode(req.UserID)
			if err != nil {
				if err.Error() == "a recent verification code already exists and is still valid" {
					// If the specific error is about an existing verification code
					res.ResponseError(c, http.StatusConflict, types.VerificationCodeExists())
					return
				}
				// Handle other internal errors
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// send the verification code to the user
			err = verificationUtils.SendRegisterVerificationCode(req.UserName, req.Email, verificationCode)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			res.ResponseSuccess(c, http.StatusOK, "request-email-verification", types.Success())

		case types.ResetPasswordEvent:
			// generate a verification code
			verificationCode, err := verificationUtils.GenerateResetPasswordVerificationCode(req.UserID)
			if err != nil {
				if err.Error() == "a recent verification code already exists and is still valid" {
					// If the specific error is about an existing verification code
					res.ResponseError(c, http.StatusConflict, types.VerificationCodeExists())
					return
				}
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// send the verification code to the user
			err = verificationUtils.SendResetPasswordVerificationCode(req.UserName, req.Email, verificationCode)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// generate a session for the user
			err = verificationUtils.GenerateVerifiedSession(ctx, req.UserID, types.ResetPasswordEvent)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// return verification success
			res.ResponseSuccess(c, http.StatusOK, "request-email-verification", types.Success())

		default:
			log.Println("Invalid event: ", req.Event)
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}
	}
}

// func VerifyEmailVerificationHandler - verifiy the verification code
func VerifyEmailVerificationHandler(verificationUtils utils.IVerificationUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req types.EmailVerifyRequest
		if err := c.BindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		user, err := verificationUtils.GetUserInfo(c)
		if err != nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		// idnetify the event
		switch req.Event {
		case types.RegisterEvent:
			// verify the verification code
			// query the latest verification code for the email
			latestVerificationCode, err := verificationUtils.GetLatestRegisterVerificationCode(user.UserID)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// check if the verification code is correct
			if latestVerificationCode != req.VerificationCode {
				res.ResponseError(c, http.StatusBadRequest, types.InvalidVerification())
				return
			}

			// hit user_server to set the user's email to verified
			err = verificationUtils.RequestEmailVerified(req.Email)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			res.ResponseSuccess(c, http.StatusOK, "verify-email", types.Success())

		case types.ResetPasswordEvent:
			// verify the verification code
			// query the latest verification code for the email
			latestVerificationCode, err := verificationUtils.GetLatestResetPasswordVerificationCode(user.UserID)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// check if the verification code is correct
			if latestVerificationCode != req.VerificationCode {
				res.ResponseError(c, http.StatusBadRequest, types.InvalidVerification())
				return
			}

			res.ResponseSuccess(c, http.StatusOK, "verify-email", types.Success())

		default:
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}
	}
}
