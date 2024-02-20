package utils

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"server/db"
	"server/schema"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

// IUserUtils is the interface for the user utils for mocking
type IUserUtils interface {
	GenerateRegisterVerificationCode(userID uint) (string, error)
	SendRegisterVerificationCode(user schema.User, code string) error
	GetUserByEmail(email string) (schema.User, error)
	CreateUser(username, email, hashedPassword string) (schema.User, error)
	ValidatePassword(password string) error
	HashPassword(password string) (string, error)
}

type UserUtils struct {
	DB db.Database
}

// Ensure UserUtils implements IUserUtils
var _ IUserUtils = (*UserUtils)(nil)

func NewUserUtils(db db.Database) *UserUtils {
	return &UserUtils{DB: db}
}

// generateRegisterVerificationCode generates a random 7-digit code for email verification, and stores it in the database
func (u *UserUtils) GenerateRegisterVerificationCode(userID uint) (string, error) {
	var registerVerification schema.RegisterEmailVerification

	// set the user id
	registerVerification.UserID = userID

	// generate a random 6-digit code
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	verificationCode := fmt.Sprintf("%07d", rng.Intn(1000000))
	registerVerification.VerificationCode = verificationCode

	// set the expiration time to 5 minutes from now
	registerVerification.ExpirationTime = time.Now().Add(5 * time.Minute)

	// create the verification record in the database everytime (for monitoring purposes), return an error if it fails
	return verificationCode, u.DB.Create(&registerVerification).Error
}

// sendRegisterVerificationCode sends an email to the user with the verification code
func (u *UserUtils) SendRegisterVerificationCode(user schema.User, code string) error {
	// email info
	fromName := os.Getenv("FROM_NAME")
	fromEmail := os.Getenv("FROM_EMAIL")
	subject := "Verification Code"
	toName := user.UserName
	toEmail := user.Email
	plainTextContent := ""
	htmlContent := "Your email verification code is " + "<strong>" + code + "</strong>.<br><br>" + "Please verify in 5 minutes."
	err := SendEmail(fromName, fromEmail, subject, toName, toEmail, plainTextContent, htmlContent)
	if err != nil {
		return errors.New("send verification email fail")
	}

	return nil
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

// GetUserByEmail retrieves a user by email
func (u *UserUtils) GetUserByEmail(email string) (schema.User, error) {
	var user schema.User
	err := u.DB.Where("email = ?", email).First(&user).Error

	// if the user is not found, return an empty user
	if err != nil {
		return schema.User{}, err
	}

	return user, nil
}

// CreateUser creates a user
func (u *UserUtils) CreateUser(username, email, hashedPassword string) (schema.User, error) {
	// create the user
	user := schema.User{
		UserName:       username,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	err := u.DB.Create(&user).Error
	if err != nil {
		return schema.User{}, err
	}

	return user, nil
}

// ValidatePassword checks if a password is valid (not empty, and at least 8 characters, includes a number, and includes a special character)
func (u *UserUtils) ValidatePassword(password string) error {
	// check if the password is at least 8 characters
	if len(password) < 8 {
		return errors.New("password is less than 8 characters")
	}

	// check if the password includes a number
	hasNumber := false
	hasCapitalLetter := false
	for _, char := range password {
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
		if char >= 'A' && char <= 'Z' {
			hasCapitalLetter = true
		}

		if hasNumber && hasCapitalLetter {
			break
		}
	}

	if !hasNumber {
		return errors.New("password does not include a number")
	}
	if !hasCapitalLetter {
		return errors.New("password does not include a capital letter")
	}

	// check if the password includes a special character (one of ?!$%^&*_+-=<>?)
	hasSpecialCharacter := false
	specialCharacters := "?!$%^&*_+-=<>?"
	for _, char := range password {
		if strings.ContainsRune(specialCharacters, char) {
			hasSpecialCharacter = true
			break
		}
	}
	if !hasSpecialCharacter {
		return errors.New("password does not include a special character")
	}

	return nil
}

// HashPassword hashes a password using bcrypt
func (u *UserUtils) HashPassword(password string) (string, error) {
	// generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
