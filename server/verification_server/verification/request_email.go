package verification

import (
	"net/http"
	"regexp"
	"verification_server/utils"

	"github.com/gin-gonic/gin"
)

// GetEmailVerificationRequest
type GetEmailVerificationRequest struct {
	Event    string `json:"event" binding:"required"`
	UserID   uint   `json:"userID" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// func RequestEmailVerificationHandler - verifiy the verification code
func RequestEmailVerificationHandler(verificationUtils utils.IVerificationUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req GetEmailVerificationRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, GeneralVerificationResponse{
				Code: "40001",
				Msg:  "Invalid request body",
			})
			return
		}

		// Check if the email is in the correct format
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+@purdue\.edu$`, req.Email)
		if !matched {
			c.JSON(http.StatusBadRequest, GeneralVerificationResponse{
				Code: "40002",
				Msg:  "Email must be a @purdue.edu address",
			})
			return
		}

		// idnetify the event
		switch req.Event {
		case RegisterEvent:
			// generate a verification code
			verificationCode, err := verificationUtils.GenerateRegisterVerificationCode(req.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, GeneralVerificationResponse{
					Code: "50001",
					Msg:  "Internal server error",
				})
				return
			}

			// send the verification code to the user
			err = verificationUtils.SendRegisterVerificationCode(req.UserName, req.Email, verificationCode)
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
			// generate a verification code
			verificationCode, err := verificationUtils.GenerateResetPasswordVerificationCode(req.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, GeneralVerificationResponse{
					Code: "50001",
					Msg:  "Internal server error",
				})
				return
			}

			// send the verification code to the user
			err = verificationUtils.SendResetPasswordVerificationCode(req.UserName, req.Email, verificationCode)
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

		default:
			c.JSON(http.StatusBadRequest, GeneralVerificationResponse{
				Code: "40004",
				Msg:  "Invalid event",
			})
			return
		}
	}
}
