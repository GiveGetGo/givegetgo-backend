package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"verification_server/schema"

	"github.com/redis/go-redis/v9"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gorm.io/gorm"
)

type VerificationUtils struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

type IVerificationUtils interface {
	GenerateRegisterVerificationCode(userID uint) (string, error)
	SendRegisterVerificationCode(username string, email string, code string) error
	GetLatestRegisterVerificationCode(userID uint) (string, error)
	RequestEmailVerified(email string) error
	GenerateResetPasswordVerificationCode(userID uint) (string, error)
	SendResetPasswordVerificationCode(username string, email string, code string) error
	GetLatestResetPasswordVerificationCode(userID uint) (string, error)
	GenerateVerifiedSession(ctx context.Context, userID uint, event string) error
}

func NewVerificationUtils(db *gorm.DB, redisClient *redis.Client) *VerificationUtils {
	return &VerificationUtils{DB: db, RedisClient: redisClient}
}

// generateRegisterVerificationCode generates a random 7-digit code for email verification, and stores it in the database
func (u *VerificationUtils) GenerateRegisterVerificationCode(userID uint) (string, error) {
	var registerVerification schema.RegisterEmailVerification

	// set the user id
	registerVerification.UserID = userID

	// generate a random 7-digit code
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	verificationCode := fmt.Sprintf("%07d", rng.Intn(1000000))
	registerVerification.VerificationCode = verificationCode

	// set the expiration time to 5 minutes from now
	registerVerification.ExpirationTime = time.Now().Add(5 * time.Minute)

	// create the verification record in the database everytime (for monitoring purposes), return an error if it fails
	return verificationCode, u.DB.Create(&registerVerification).Error
}

// sendRegisterVerificationCode sends an email to the user with the verification code
func (u *VerificationUtils) SendRegisterVerificationCode(username string, email string, code string) error {
	// email info
	fromName := os.Getenv("FROM_NAME")
	fromEmail := os.Getenv("FROM_EMAIL")
	subject := "Registeration Verification Code"
	toName := username
	toEmail := email
	plainTextContent := ""
	htmlContent := "Your email verification code is " + "<strong>" + code + "</strong>.<br><br>" + "Please verify in 5 minutes."
	err := SendEmail(fromName, fromEmail, subject, toName, toEmail, plainTextContent, htmlContent)
	if err != nil {
		return errors.New("send verification email fail")
	}

	return nil
}

// GenerateResetPasswordVerificationCode generates a random 7-digit code for reset password, and stores it in the database
func (u *VerificationUtils) GenerateResetPasswordVerificationCode(userID uint) (string, error) {
	var resetPasswordVerification schema.ResetPasswordVerification

	// set the user id
	resetPasswordVerification.UserID = userID

	// generate a random 7-digit code
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	resetCode := fmt.Sprintf("%07d", rng.Intn(1000000))
	resetPasswordVerification.ResetCode = resetCode

	// set the expiration time to 5 minutes from now
	resetPasswordVerification.ExpirationTime = time.Now().Add(5 * time.Minute)

	// create the verification record in the database everytime (for monitoring purposes), return an error if it fails
	return resetCode, u.DB.Create(&resetPasswordVerification).Error
}

// RequestEmailVerified sets the user's email to verified
func (u *VerificationUtils) RequestEmailVerified(email string) error {
	// hit user_server to set the user's email to verified
	requestBody, err := json.Marshal(struct {
		Email string `json:"email"`
	}{
		Email: email,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("USER_SERVICE_URL")+"/v1/internal/user/email-verified", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
		return err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service", "VERIFICATION")                    // Set the service name
	req.Header.Set("X-Api-Key", os.Getenv("VERIFICATION_API_KEY")) // Set the API key

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("user service responded with status: ", resp.StatusCode)
		// Handle non-OK responses here
		return fmt.Errorf("user service responded with status: %d", resp.StatusCode)
	}

	log.Println("successfully set user email to verified")

	return nil
}

// SendResetPasswordVerificationCode sends an email to the user with the verification code
func (u *VerificationUtils) SendResetPasswordVerificationCode(username string, email string, code string) error {
	// email info
	fromName := os.Getenv("FROM_NAME")
	fromEmail := os.Getenv("FROM_EMAIL")
	subject := "Reset Password Verification Code"
	toName := username
	toEmail := email
	plainTextContent := ""
	htmlContent := "Your reset password verification code is " + "<strong>" + code + "</strong>.<br><br>" + "Please reset in 5 minutes."
	err := SendEmail(fromName, fromEmail, subject, toName, toEmail, plainTextContent, htmlContent)
	if err != nil {
		return errors.New("send verification email fail")
	}

	return nil
}

// GetLatestResetPasswordVerificationCode returns the latest verification code for the user
func (u *VerificationUtils) GetLatestResetPasswordVerificationCode(userID uint) (string, error) {
	var resetPasswordVerification schema.ResetPasswordVerification
	// get the latest verification code with the user id
	err := u.DB.Where("user_id = ?", userID).Order("created_at desc").First(&resetPasswordVerification).Error
	if err != nil {
		return "", err
	}

	return resetPasswordVerification.ResetCode, nil
}

// send email func
func SendEmail(fromName, fromEmail, subject, toName, toEmail, plainTextContent, htmlContent string) error {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// func GetLatestRegisterVerificationCode returns the latest verification code for the user
func (u *VerificationUtils) GetLatestRegisterVerificationCode(userID uint) (string, error) {
	var registerVerification schema.RegisterEmailVerification
	// get the latest verification code with the user id
	err := u.DB.Where("user_id = ?", userID).Order("created_at desc").First(&registerVerification).Error
	if err != nil {
		return "", err
	}

	return registerVerification.VerificationCode, nil
}

// Set redis session after verification
func (u *VerificationUtils) GenerateVerifiedSession(ctx context.Context, userID uint, event string) error {
	sessionKey := fmt.Sprintf("session:%d:%s", userID, event) // Set the session key
	sessionValue := "verified"                                // Set the session value
	expiration := 5 * time.Minute                             // Set the expiration time for the session

	// Set the session key, value, and expiration in Redis
	err := u.RedisClient.Set(ctx, sessionKey, sessionValue, expiration).Err()
	if err != nil {
		return err
	}

	// print the session key
	log.Println("session key:", sessionKey)

	return nil // Return nil if no errors occurred
}
