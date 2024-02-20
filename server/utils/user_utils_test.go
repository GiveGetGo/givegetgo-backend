package utils

import (
	"log"
	"server/db"
	schema "server/schema"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGenerateRegisterVerificationCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock database object
	mockDB := db.NewMockDatabase(ctrl)
	userUtils := NewUserUtils(mockDB)

	log.Println("TestGenerateRegisterVerificationCode")
	log.Println(mockDB)
	log.Println(userUtils)

	// Test case for successful code generation
	t.Run("successful code generation", func(t *testing.T) {
		userID := uint(1)
		mockDB.EXPECT().Create(gomock.Any()).Return(&gorm.DB{Error: nil}) // Expect database create call to succeed

		code, err := userUtils.GenerateRegisterVerificationCode(userID)
		assert.NoError(t, err)
		assert.Len(t, code, 7) // Ensure code is 7 digits
	})
}

func TestSendRegisterVerificationCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock database object
	mockDB := db.NewMockDatabase(ctrl)
	userUtils := NewUserUtils(mockDB)

	// Test case for successful code generation
	t.Run("successful code generation", func(t *testing.T) {
		user := schema.User{
			UserName: "testuser",
			Email:    "root@givegetgo.xyz",
		}
		code := "1234567"

		err := userUtils.SendRegisterVerificationCode(user, code)
		assert.NoError(t, err)
	})
}

func TestGetUserByEmail(t *testing.T) {
	// Create a new mock SQL database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	// Set up GORM to use the mock database
	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open mock GORM database")
	}

	userUtils := NewUserUtils(db)

	email := "test@purdue.edu"
	expectedUser := schema.User{UserID: 1, UserName: "testuser", Email: email}

	// Test case for user found
	t.Run("user found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"UserID", "UserName", "Email"}).
			AddRow(1, "testuser", "test@purdue.edu")
		mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

		foundUser, err := userUtils.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, foundUser)
	})

	// Test case for user not found
	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT`).WillReturnError(gorm.ErrRecordNotFound)

		foundUser, err := userUtils.GetUserByEmail(email)
		assert.Error(t, err)
		assert.Equal(t, schema.User{}, foundUser)
	})

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock database object
	mockDB := db.NewMockDatabase(ctrl)
	userUtils := NewUserUtils(mockDB)

	// Test case for user creation success
	t.Run("Created User", func(t *testing.T) {
		userName := "testuser"
		email := "test@purdue.edu"
		hashedPassword := "hashedPassword"
		mockDB.EXPECT().Create(&schema.User{UserName: userName, Email: email, HashedPassword: hashedPassword}).Return(&gorm.DB{Error: nil}) // Expect database create call to succeed

		user, err := userUtils.CreateUser(userName, email, hashedPassword)
		assert.NoError(t, err)
		assert.NotEqual(t, user, schema.User{})
	})
}

func TestValidatePassword(t *testing.T) {
	userUtils := NewUserUtils(nil)

	t.Run("Valid Password", func(t *testing.T) {
		err := userUtils.ValidatePassword("Password123!")
		assert.NoError(t, err)
	})

	t.Run("Invalid Password, too short", func(t *testing.T) {
		err := userUtils.ValidatePassword("Pass1!")
		assert.Error(t, err)
	})

	t.Run("Invalid Password, no uppercase", func(t *testing.T) {
		err := userUtils.ValidatePassword("password123!")
		assert.Error(t, err)
	})

	t.Run("Invalid Password, no number", func(t *testing.T) {
		err := userUtils.ValidatePassword("password")
		assert.Error(t, err)
	})
}

func TestHashPassword(t *testing.T) {
	userUtils := NewUserUtils(nil)

	hashedPassword, err := userUtils.HashPassword("password123")
	assert.NoError(t, err)

	// Verify the password is correctly hashed
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte("password123"))
	assert.NoError(t, err)
}
