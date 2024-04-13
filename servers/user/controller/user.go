package controller

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"user/schema"
	"user/utils"

	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterHandler is the handler for the user registration endpoint
func RegisterHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		log.Println("RegisterHandler called")

		var req types.RegisterRequest
		//if err := c.BindJSON(&req); err != nil {
		//	types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
		//	return
		//}

		if err := c.BindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Check if the email is in the correct format
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+@purdue\.edu$`, req.Email)
		if !matched {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidEmail())
			return
		}

		// Check if the email is already registered
		user, err := userUtils.GetUserByEmail(req.Email)

		// If no error, or the error is not a record not found error, return an internal server error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// If the user exists, return an error
		if (schema.User{}) != user {
			types.ResponseError(c, http.StatusBadRequest, types.AlreadyExists())
			return
		}

		// Check if the password is valid
		err = userUtils.ValidatePassword(req.Password)
		if err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidCredentials())
			return
		}

		// hash the password
		hashedPassword, err := userUtils.HashPassword(req.Password)
		if err != nil {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Create the user
		user, err = userUtils.CreateUser(req.UserName, req.Email, hashedPassword)
		if err != nil {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// request the verification server to send a verification email
		err = userUtils.RequestRegisterVerificationEmail(user.UserID, req.UserName, req.Email)
		if err != nil {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return the user
		types.ResponseSuccess(c, http.StatusCreated, "register", user.UserID, types.UserCreated())
	}
}

func LoginHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Check if the email is already registered
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// If email not verified, return an error
		if !user.EmailVerified {
			types.ResponseError(c, http.StatusBadRequest, types.EmailNotVerified())
			return
		}

		authenticated := userUtils.AuthenticateUser(user, req.Password)
		if !authenticated {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidCredentials())
			return
		}

		types.ResponseSuccess(c, http.StatusOK, "login", user.UserID, types.LoginSuccess())
	}
}

func ForgotPasswordHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserForgetPassRequest
		// Parse and validate the request body
		if err := c.ShouldBindJSON(&req); err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Check if the email is in the correct format
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+@purdue\.edu$`, req.Email)
		if !matched {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidEmail())
			return
		}

		// Check if the email exists in your database
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				types.ResponseError(c, http.StatusBadRequest, types.UserNotFound())
				return
			}

			// Handle internal server error
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// call /verification/request-email endpoint in the verification server to send the verification email
		// if the verification server is down, return an internal server error
		err = userUtils.RequestForgetpassVerificationEmail(user.UserID, user.UserName, req.Email)
		if err != nil {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		types.ResponseSuccess(c, http.StatusOK, "forgot-password", user.UserID, types.Success())
	}
}

// ResetPasswordHandler is the handler for the reset password endpoint
func ResetPasswordHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req types.UserResetPassRequest
		if err := c.BindJSON(&req); err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// find the user by email
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				types.ResponseError(c, http.StatusBadRequest, types.UserNotFound())
				return
			}

			// Handle internal server error
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		ctx := c.Request.Context() // get the context

		// check if session is valid
		err = userUtils.CheckEmailVerificationSession(ctx, user.UserID, types.ResetPasswordEvent)
		if err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidSession())
			return
		}

		// Check if the password is valid
		err = userUtils.ValidatePassword(req.Newpassword)
		if err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidCredentials())
			return
		}

		// hash the password
		hashedPassword, err := userUtils.HashPassword(req.Newpassword)
		if err != nil {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Update the user's password
		err = userUtils.UpdatePassword(user.UserID, hashedPassword)
		if err != nil {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return success
		types.ResponseSuccess(c, http.StatusOK, "reset-password", user.UserID, types.Success())
	}
}

// successful email verification from verification_server, call this endpoint to set the user's email to verified
func SetUserEmailVerifiedHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req types.SetUserEmailVerifiedRequest
		if err := c.BindJSON(&req); err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Update the user's email to verified
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			log.Println(err)
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// If the user not found, return an error
		if (schema.User{}) == user {
			types.ResponseError(c, http.StatusBadRequest, types.UserNotFound())
			return
		}

		// If the user's email is already verified, return an error
		if user.EmailVerified {
			types.ResponseError(c, http.StatusBadRequest, types.AlreadyVerified())
			return
		}

		// Update the user's email to verified
		err = userUtils.MarkEmailVerified(req.Email)
		if err != nil {
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return success
		types.ResponseSuccess(c, http.StatusOK, "set-email-verified", user.UserID, types.EmailVerified())
	}
}
