package controller

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"user/schema"
	"user/utils"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterHandler is the handler for the user registration endpoint
func RegisterHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req types.RegisterRequest
		if err := c.BindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Check if the email is in the correct format
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+@purdue\.edu$`, req.Email)
		if !matched {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidEmail())
			return
		}

		// Check if the email is already registered
		user, err := userUtils.GetUserByEmail(req.Email)

		// If no error, or the error is not a record not found error, return an internal server error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// If the user exists, return an error
		if (schema.User{}) != user {
			res.ResponseError(c, http.StatusBadRequest, types.AlreadyExists())
			return
		}

		// Check if the password is valid
		err = userUtils.ValidatePassword(req.Password)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidCredentials())
			return
		}

		// hash the password
		hashedPassword, err := userUtils.HashPassword(req.Password)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Extract the username from the email
		splitEmail := strings.Split(req.Email, "@")
		if len(splitEmail) == 0 {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidEmail())
			return
		}
		username := splitEmail[0]

		// Create the user with the extracted username
		user, err = userUtils.CreateUser(username, req.Email, hashedPassword, req.Class, req.Major)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// request the verification server to send a verification email
		err = userUtils.RequestRegisterVerificationEmail(user.UserID, req.Email, req.Email)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// add user id to session
		session := sessions.Default(c)
		session.Set("userid", user.UserID)
		err = session.Save()
		if err != nil {
			log.Println("err is ", err)
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return the user
		res.ResponseSuccess(c, http.StatusCreated, "register", types.UserCreated())
	}
}

func LoginHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Check if the email is already registered
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// If email not verified, return an error
		if !user.EmailVerified {
			// resend verification email if not verified
			err = userUtils.RequestRegisterVerificationEmail(user.UserID, req.Email, req.Email)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}

			res.ResponseError(c, http.StatusBadRequest, types.EmailNotVerified())
			return
		}

		authenticated := userUtils.AuthenticateUser(user, req.Password)
		if !authenticated {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidCredentials())
			return
		}

		// set session
		session := sessions.Default(c)
		session.Set("userid", user.UserID)
		err = session.Save()
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "login", types.LoginSuccess())
	}
}

func ForgotPasswordHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.UserForgetPassRequest
		// Parse and validate the request body
		if err := c.ShouldBindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Check if the email is in the correct format
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+@purdue\.edu$`, req.Email)
		if !matched {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidEmail())
			return
		}

		// Check if the email exists in your database
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.ResponseError(c, http.StatusBadRequest, types.UserNotFound())
				return
			}

			// Handle internal server error
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// call /verification/request-email endpoint in the verification server to send the verification email
		// if the verification server is down, return an internal server error
		err = userUtils.RequestForgetpassVerificationEmail(user.UserID, user.UserName, req.Email)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "forgot-password", types.Success())
	}
}

// ResetPasswordHandler is the handler for the reset password endpoint
func ResetPasswordHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req types.UserResetPassRequest
		if err := c.BindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// find the user by email
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.ResponseError(c, http.StatusBadRequest, types.UserNotFound())
				return
			}

			// Handle internal server error
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		ctx := c.Request.Context() // get the context

		// check if session is valid
		err = userUtils.CheckEmailVerificationSession(ctx, user.UserID, types.ResetPasswordEvent)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidSession())
			return
		}

		// Check if the password is valid
		err = userUtils.ValidatePassword(req.Newpassword)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidCredentials())
			return
		}

		// hash the password
		hashedPassword, err := userUtils.HashPassword(req.Newpassword)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Update the user's password
		err = userUtils.UpdatePassword(user.UserID, hashedPassword)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return success
		res.ResponseSuccess(c, http.StatusOK, "reset-password", types.Success())
	}
}

// successful email verification from verification_server, call this endpoint to set the user's email to verified
func SetUserEmailVerifiedHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var req types.SetUserEmailVerifiedRequest
		if err := c.BindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Update the user's email to verified
		user, err := userUtils.GetUserByEmail(req.Email)
		if err != nil {
			log.Println(err)
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// If the user not found, return an error
		if (schema.User{}) == user {
			res.ResponseError(c, http.StatusBadRequest, types.UserNotFound())
			return
		}

		// If the user's email is already verified, return an error
		if user.EmailVerified {
			res.ResponseError(c, http.StatusBadRequest, types.AlreadyVerified())
			return
		}

		// Update the user's email to verified
		err = userUtils.MarkEmailVerified(req.Email)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return success
		res.ResponseSuccess(c, http.StatusOK, "set-email-verified", types.EmailVerified())
	}
}

// Logout handler for session termination
func LogoutHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		err := session.Save()
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
		}

		res.ResponseSuccess(c, http.StatusOK, "logout", types.Success())
	}
}

// Handle user session for internal session authentication between services
func SessionHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		res.ResponseSuccess(c, http.StatusOK, "session", types.Success())
	}
}

// User Info
func GetMeHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
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
			res.ResponseError(c, http.StatusNotFound, types.UserNotFound())
			return
		}

		responseInfo := types.UserInfoResponse{
			UserID:        user.UserID,
			Username:      user.UserName,
			Email:         user.Email,
			Class:         user.Class,
			Major:         user.Major,
			EmailVerified: user.EmailVerified,
			MfaVerified:   user.MFAVerified,
		}

		res.ResponseSuccessWithData(c, http.StatusOK, "get user info", types.Success(), responseInfo)
	}
}

func EditMeHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userIdValue := session.Get("userid")
		if userIdValue == nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		userId, ok := userIdValue.(uint)
		if !ok {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		var updateReq types.UserUpdateRequest
		if err := c.BindJSON(&updateReq); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		if err := userUtils.UpdateUser(userId, updateReq); err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "User updated", types.Success())
	}
}

func VerifiedHandler(userUtils utils.IUserUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userIdValue := session.Get("userid")
		if userIdValue == nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		userId, ok := userIdValue.(uint)
		if !ok {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		user, err := userUtils.GetUserByID(userId)
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

		res.ResponseSuccess(c, http.StatusOK, "verified", types.Success())
	}
}
