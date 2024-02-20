package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"user_server/schema"
	"user_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestUserRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock UserUtils
	mockUserUtils := utils.NewMockIUserUtils(ctrl)

	// Define your test cases
	testCases := []struct {
		name           string
		requestBody    UserRegisterRequest
		mockSetup      func()
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "Successful Registration",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUserUtils.EXPECT().GetUserByEmail("test@purdue.edu").Return(schema.User{}, gorm.ErrRecordNotFound)
				mockUserUtils.EXPECT().ValidatePassword("Password123!").Return(nil)
				mockUserUtils.EXPECT().HashPassword("Password123!").Return("hashedPassword", nil)
				mockUserUtils.EXPECT().CreateUser("testuser", "test@purdue.edu", "hashedPassword").Return(schema.User{UserID: 1}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "20101",
		},
		{
			name: "Invalid Request Body",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				// Missing the password field
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "40001",
		},
		{
			name: "Invalid Email",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@gmail.com",
				Password: "password123",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "40002",
		},
		{
			name: "Email Already Registered",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				Password: "password123",
			},
			mockSetup: func() {
				// Simulate that the user already exists
				existingUser := schema.User{
					UserID:         1,
					UserName:       "testuser",
					Email:          "test@purdue.edu",
					HashedPassword: "hashedPassword",
				}
				mockUserUtils.EXPECT().GetUserByEmail("test@purdue.edu").Return(existingUser, nil)
				// Note that we don't set up expectations for CreateUser, GenerateRegisterVerificationCode, or SendRegisterVerificationCode
				// as these should not be called when a user already exists
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "40003",
		},
		{
			name: "Invalid Password, less than 8 characters",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				Password: "pass123", // Less than 8 characters
			},
			mockSetup: func() {
				mockUserUtils.EXPECT().GetUserByEmail("test@purdue.edu").Return(schema.User{}, gorm.ErrRecordNotFound)
				mockUserUtils.EXPECT().ValidatePassword("pass123").Return(errors.New("invalid password, less than 8 characters"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "40004",
		},
		{
			name: "Invalid Password, does not include a number",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				Password: "password", // Does not include a number
			},
			mockSetup: func() {
				mockUserUtils.EXPECT().GetUserByEmail("test@purdue.edu").Return(schema.User{}, gorm.ErrRecordNotFound)
				mockUserUtils.EXPECT().ValidatePassword("password").Return(errors.New("invalid password, does not include a number"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "40004",
		},
		{
			name: "Invalid Password, does not include a capital letter",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				Password: "password123", // Does not include a capital letter
			},
			mockSetup: func() {
				mockUserUtils.EXPECT().GetUserByEmail("test@purdue.edu").Return(schema.User{}, gorm.ErrRecordNotFound)
				mockUserUtils.EXPECT().ValidatePassword("password123").Return(errors.New("invalid password, does not include a capital letter"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "40004",
		},
		{
			name: "Invalid Password, does not include a special character",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				Password: "Password123", // Does not include a special character
			},
			mockSetup: func() {
				mockUserUtils.EXPECT().GetUserByEmail("test@purdue.edu").Return(schema.User{}, gorm.ErrRecordNotFound)
				mockUserUtils.EXPECT().ValidatePassword("Password123").Return(errors.New("invalid password, does not include a special character"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "40004",
		},
		{
			name: "Failed to Hash Password",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUserUtils.EXPECT().GetUserByEmail("test@purdue.edu").Return(schema.User{}, gorm.ErrRecordNotFound)
				mockUserUtils.EXPECT().ValidatePassword("Password123!").Return(nil)
				mockUserUtils.EXPECT().HashPassword("Password123!").Return("", errors.New("hash password error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "50002",
		},
		{
			name: "Failed to Create User",
			requestBody: UserRegisterRequest{
				UserName: "testuser",
				Email:    "test@purdue.edu",
				Password: "Password123!",
			},
			mockSetup: func() {
				mockUserUtils.EXPECT().GetUserByEmail("test@purdue.edu").Return(schema.User{}, gorm.ErrRecordNotFound)
				mockUserUtils.EXPECT().ValidatePassword("Password123!").Return(nil)
				mockUserUtils.EXPECT().HashPassword("Password123!").Return("hashedPassword", nil)
				mockUserUtils.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(schema.User{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "50003",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup the mock expectations
			tc.mockSetup()

			// Marshal the request body to JSON
			body, _ := json.Marshal(tc.requestBody)

			// Create a request to the handler
			w := httptest.NewRecorder()
			ctx, r := gin.CreateTestContext(w)
			r.POST("/user/register", UserRegisterHandler(mockUserUtils))
			req, _ := http.NewRequest(http.MethodPost, "/user/register", bytes.NewBuffer(body))
			ctx.Request = req

			// Call the handler
			r.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tc.expectedStatus, w.Code)
			var response UserRegisterResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, response.Code)
		})
	}
}
