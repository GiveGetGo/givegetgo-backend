package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"user/db"
	"user/middleware"
	"user/schema"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
)

// IUserUtils is the interface for the user utils for mocking
type IUserUtils interface {
	// Create
	CreateUser(username, email, hashedPassword string, class string, major string) (schema.User, error)

	// Get info
	GetUserByID(userID uint) (schema.User, error)
	GetUserByEmail(email string) (schema.User, error)

	// Update
	UpdateUser(userID uint, updates types.UserUpdateRequest) error
	UpdatePassword(userID uint, hashedPassword string) error

	// Delete
	DeleteUser(userID uint) error

	// Others
	ValidatePassword(password string) error
	HashPassword(password string) (string, error)
	AuthenticateUser(user schema.User, password string) bool
	RequestRegisterVerificationEmail(userID uint, username string, email string) error
	RequestForgetpassVerificationEmail(userID uint, username string, email string) error
	MarkEmailVerified(email string) error
	MarkMFAVerified(userID uint) error
	StoreEncryptedTOTPSecret(userID uint, encryptedSecret string) error
	CheckEmailVerificationSession(ctx context.Context, userID uint, event string) error
	GenerateAndSendQRCode(c *gin.Context, email string, secret []byte)
}

type UserUtils struct {
	DB          db.Database
	RedisClient middleware.RedisClientInterface
}

// Ensure UserUtils implements IUserUtils
var _ IUserUtils = (*UserUtils)(nil)

func NewUserUtils(db db.Database, redisClient middleware.RedisClientInterface) *UserUtils {
	return &UserUtils{DB: db, RedisClient: redisClient}
}

// GetUserByID retrieves a user by ID
func (u *UserUtils) GetUserByID(userID uint) (schema.User, error) {
	var user schema.User
	err := u.DB.Where("user_id = ?", userID).First(&user).Error

	// if the user is not found, return an empty user
	if err != nil {
		return schema.User{}, err
	}

	return user, nil
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
func (u *UserUtils) CreateUser(username, email, hashedPassword string, class string, major string) (schema.User, error) {
	// create the user
	user := schema.User{
		UserName:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		Class:          class,
		Major:          major,
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

func (u *UserUtils) AuthenticateUser(user schema.User, password string) bool {
	// Compare the hash of the input password with the hash stored in the database.
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	return err == nil
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
	req, err := http.NewRequest("POST", os.Getenv("VERIFICATION_SERVICE_URL")+"/v1/internal/verification/request-email", bytes.NewBuffer(verificationReqBody))
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

	if resp.StatusCode == http.StatusConflict {
		log.Println("A recent verification code already exists, no new code sent.")
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("verification service responded with status: ", resp.StatusCode)
		// Handle non-OK responses here
		return fmt.Errorf("verification service responded with status: %d", resp.StatusCode)
	}

	log.Println("successfully sent verification code")

	return nil
}

// RequestVerificationEmail - request verification email through calling verification_server /verification/request-email
func (u *UserUtils) RequestForgetpassVerificationEmail(userID uint, username string, email string) error {
	verificationReqBody, err := json.Marshal(struct {
		Event    string `json:"event"`
		UserID   uint   `json:"userID"`
		UserName string `json:"username"`
		Email    string `json:"email"`
	}{
		Event:    "reset-password",
		UserID:   userID,
		UserName: username,
		Email:    email,
	})

	if err != nil {
		log.Println("error marshalling request body for forget password")
		return err
	}

	// Create the HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("POST", os.Getenv("VERIFICATION_SERVICE_URL")+"/v1/internal/verification/request-email", bytes.NewBuffer(verificationReqBody))
	if err != nil {
		log.Println("error creating request for forget password:", err)
		return err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service", "USER")                    // Set the service name
	req.Header.Set("X-Api-Key", os.Getenv("USER_API_KEY")) // Set the API key

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error sending forget password verification email:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		log.Println("A recent verification code already exists, no new code sent.")
		return nil
	}

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		log.Println("verification service responded with status for forget password:", resp.StatusCode)
		return fmt.Errorf("verification service responded with status: %d", resp.StatusCode)
	}

	log.Println("successfully sent forget password verification email")
	return nil
}

// MarkEmailVerified - set user email verified
func (u *UserUtils) MarkEmailVerified(email string) error {
	// update the user's email to verified
	err := u.DB.Model(&schema.User{}).Where("email = ?", email).Update("email_verified", true).Error
	if err != nil {
		return err
	}

	return nil
}

// MarkMFAVerified - mark MFA verified
func (u *UserUtils) MarkMFAVerified(userID uint) error {
	// mark the user as MFA verified
	err := u.DB.Model(&schema.User{}).Where("user_id = ?", userID).Update("mfa_verified", true).Error
	if err != nil {
		return err
	}

	return nil
}

// StoreEncryptedTOTPSecret - store encrypted TOTP secret
func (u *UserUtils) StoreEncryptedTOTPSecret(userID uint, encryptedSecret string) error {
	// store the encrypted secret in the database
	err := u.DB.Model(&schema.User{}).Where("user_id = ?", userID).Update("mfa_secret", encryptedSecret).Error
	if err != nil {
		return err
	}

	return nil
}

// CheckEmailVerificationSession checks if the user's email verification session exists and is valid
func (u *UserUtils) CheckEmailVerificationSession(ctx context.Context, userID uint, event string) error {
	sessionKey := fmt.Sprintf("session:%d:%s", userID, event) // Construct the session key

	// Use the Redis GET command to retrieve the session value
	sessionValue, err := u.RedisClient.Get(ctx, sessionKey).Result()

	if err == redis.Nil {
		// The key does not exist or the session has expired
		return fmt.Errorf("session not found or expired for user ID %d and event %s", userID, event)
	} else if err != nil {
		// An error occurred while trying to retrieve the session value
		return err
	}

	// Check the session value to determine if it's the expected one
	if sessionValue != "verified" {
		// The session value is not what was expected
		return fmt.Errorf("unexpected session value for user ID %d and event %s", userID, event)
	}

	// Delete the session key from Redis and return nil
	_, err = u.RedisClient.Del(ctx, sessionKey).Result()
	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates the user's password
func (u *UserUtils) UpdatePassword(userID uint, hashedPassword string) error {
	var user schema.User
	// First, find the user by ID to ensure they exist.
	if err := u.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		log.Printf("Error finding user with ID %d: %v", userID, err)
		return err
	}

	// Update the user's hashed password
	user.HashedPassword = hashedPassword
	if err := u.DB.Save(&user).Error; err != nil {
		log.Printf("Error updating user's password: %v", err)
		return err
	}

	// Return the updated user object and nil for the error
	return nil
}

func (u *UserUtils) UpdateUser(userID uint, updates types.UserUpdateRequest) error {
	updateMap := map[string]interface{}{
		"username":      updates.Username,
		"class":         updates.Class,
		"major":         updates.Major,
		"profile_image": updates.ProfileImage,
		"profile_info":  updates.ProfileInfo,
	}
	return u.DB.Model(&schema.User{}).Where("user_id = ?", userID).Updates(updateMap).Error
}

func (u *UserUtils) DeleteUser(userID uint) error {
	result := u.DB.Delete(&schema.User{}, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no user found")
	}
	return nil
}

func (u *UserUtils) GenerateAndSendQRCode(c *gin.Context, email string, secret []byte) {
	uri := fmt.Sprintf("otpauth://totp/GiveGetGo:%s?secret=%s&issuer=GiveGetGo", email, string(secret))
	qrCode, err := qrcode.Encode(uri, qrcode.Medium, 256)
	if err != nil {
		res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
		return
	}
	c.Data(http.StatusOK, "image/png", qrCode)
}
