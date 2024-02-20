package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"user_server/db"
	"user_server/schema"

	"golang.org/x/crypto/bcrypt"
)

// IUserUtils is the interface for the user utils for mocking
type IUserUtils interface {
	GetUserByEmail(email string) (schema.User, error)
	CreateUser(username, email, hashedPassword string) (schema.User, error)
	ValidatePassword(password string) error
	HashPassword(password string) (string, error)
	RequestRegisterVerificationEmail(userID uint, username string, email string) error
	SetUserEmailVerified(email string) error
}

type UserUtils struct {
	DB db.Database
}

// Ensure UserUtils implements IUserUtils
var _ IUserUtils = (*UserUtils)(nil)

func NewUserUtils(db db.Database) *UserUtils {
	return &UserUtils{DB: db}
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

// RequestVerificationEmail - request verification email through calling verification_server /verification/request-email
func (u *UserUtils) RequestRegisterVerificationEmail(userID uint, username string, email string) error {
	verificationReqBody, err := json.Marshal(struct {
		Event    string `json:"event"`
		UserID   uint   `json:"userID"`
		UserName string `json:"username"`
		Email    string `json:"email"`
	}{
		Event:    "register",
		UserID:   userID,
		UserName: username,
		Email:    email,
	})
	if err != nil {
		log.Println("error marshalling request body")
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("VERIFICATION_SERVICE_URL")+"/v1/verification/request-email", bytes.NewBuffer(verificationReqBody))
	if err != nil {
		log.Println(err)
		return err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service", "USER")                    // Set the service name
	req.Header.Set("X-Api-Key", os.Getenv("USER_API_KEY")) // Set the API key

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("verification service responded with status: ", resp.StatusCode)
		// Handle non-OK responses here
		return fmt.Errorf("verification service responded with status: %d", resp.StatusCode)
	}

	log.Println("successfully sent verification code")

	return nil
}

// SetUserEmailVerified - set user email verified
func (u *UserUtils) SetUserEmailVerified(email string) error {
	// update the user's email to verified
	err := u.DB.Model(&schema.User{}).Where("email = ?", email).Update("email_verified", true).Error
	if err != nil {
		return err
	}

	return nil
}
