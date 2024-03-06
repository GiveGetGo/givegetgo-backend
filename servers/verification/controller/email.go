package controller

import (
	"log"
	"net/http"
	"regexp"
	"verification/utils"

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
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// get the context
		ctx := c.Request.Context()

		// Check if the email is in the correct format
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+@purdue\.edu$`, req.Email)
		if !matched {
			log.Println("Invalid email: ", req.Email)
			types.ResponseError(c, http.StatusBadRequest, types.InvalidEmail())
			return
		}

		// idnetify the event
		switch req.Event {
		case types.RegisterEvent:
			// generate a verification code
			verificationCode, err := verificationUtils.GenerateRegisterVerificationCode(req.UserID)
			if err != nil {
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// send the verification code to the user
			err = verificationUtils.SendRegisterVerificationCode(req.UserName, req.Email, verificationCode)
			if err != nil {
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			types.ResponseSuccess(c, http.StatusOK, "request-email-verification", req.UserID, types.Success())

		case types.ResetPasswordEvent:
			// generate a verification code
			verificationCode, err := verificationUtils.GenerateResetPasswordVerificationCode(req.UserID)
			if err != nil {
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// send the verification code to the user
			err = verificationUtils.SendResetPasswordVerificationCode(req.UserName, req.Email, verificationCode)
			if err != nil {
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// generate a session for the user
			err = verificationUtils.GenerateVerifiedSession(ctx, req.UserID, types.ResetPasswordEvent)
			if err != nil {
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// return verification success
			types.ResponseSuccess(c, http.StatusOK, "request-email-verification", req.UserID, types.Success())

		default:
			log.Println("Invalid event: ", req.Event)
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
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
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// idnetify the event
		switch req.Event {
		case types.RegisterEvent:
			// verify the verification code
			// query the latest verification code for the email
			latestVerificationCode, err := verificationUtils.GetLatestRegisterVerificationCode(req.UserID)
			if err != nil {
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// check if the verification code is correct
			if latestVerificationCode != req.VerificationCode {
				types.ResponseError(c, http.StatusBadRequest, types.InvalidVerification())
				return
			}

			// hit user_server to set the user's email to verified
			err = verificationUtils.RequestEmailVerified(req.Email)
			if err != nil {
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			types.ResponseSuccess(c, http.StatusOK, "verify-email", req.UserID, types.Success())

		case types.ResetPasswordEvent:
			// verify the verification code
			// query the latest verification code for the email
			latestVerificationCode, err := verificationUtils.GetLatestResetPasswordVerificationCode(req.UserID)
			if err != nil {
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			// check if the verification code is correct
			if latestVerificationCode != req.VerificationCode {
				types.ResponseError(c, http.StatusBadRequest, types.InvalidVerification())
				return
			}

			types.ResponseSuccess(c, http.StatusOK, "verify-email", req.UserID, types.Success())

		default:
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}
	}
}
