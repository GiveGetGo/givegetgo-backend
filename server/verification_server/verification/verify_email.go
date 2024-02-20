package verification

import (
	"net/http"
	"verification_server/utils"

	"github.com/gin-gonic/gin"
)

// EmailVerificationRequest
type EmailVerificationRequest struct {
	Event            string `json:"event" binding:"required"`
	UserID           uint   `json:"userID" binding:"required"`
	Email            string `json:"email" binding:"required"`
	VerificationCode string `json:"verification_code" binding:"required"`
}

// func VerifyEmailVerificationHandler - verifiy the verification code
func VerifyEmailVerificationHandler(verificationUtils utils.IVerificationUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req EmailVerificationRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, GeneralVerificationResponse{
				Code: "40001",
				Msg:  "Invalid request body",
			})
			return
		}

		// idnetify the event
		switch req.Event {
		case RegisterEvent:
			// verify the verification code
			// query the latest verification code for the email
			latestVerificationCode, err := verificationUtils.GetLatestRegisterVerificationCode(req.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, GeneralVerificationResponse{
					Code: "50001",
					Msg:  "Internal server error",
				})
				return
			}

			// check if the verification code is correct
			if latestVerificationCode != req.VerificationCode {
				c.JSON(http.StatusBadRequest, GeneralVerificationResponse{
					Code: "40003",
					Msg:  "Invalid verification code",
				})
				return
			}

			// hit user_server to set the user's email to verified
			err = verificationUtils.RequestEmailVerified(req.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, GeneralVerificationResponse{
					Code: "50001",
					Msg:  "Internal server error",
				})
				return
			}

			// return verification success
			c.JSON(http.StatusOK, GeneralVerificationResponse{
				Code: "20000",
				Msg:  "Verification success",
			})
		case ResetEvent:
			// verify the verification code
			// query the latest verification code for the email
			latestVerificationCode, err := verificationUtils.GetLatestResetPasswordVerificationCode(req.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, GeneralVerificationResponse{
					Code: "50001",
					Msg:  "Internal server error",
				})
				return
			}

			// check if the verification code is correct
			if latestVerificationCode != req.VerificationCode {
				c.JSON(http.StatusBadRequest, GeneralVerificationResponse{
					Code: "40003",
					Msg:  "Invalid verification code",
				})
				return
			}

			// return verification success
			c.JSON(http.StatusOK, GeneralVerificationResponse{
				Code: "20000",
				Msg:  "Verification success",
			})
		default:
			c.JSON(http.StatusBadRequest, GeneralVerificationResponse{
				Code: "40004",
				Msg:  "Invalid event",
			})
			return
		}
	}
}
