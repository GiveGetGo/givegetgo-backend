package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"notification/db"
	"notification/middleware"
	"notification/schema"
	"os"
	"time"

	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type INotificationUtils interface {
	GetNotificationByUserID(userID uint) ([]schema.Notification, error)
	DeleteNotificationByID(notificationID uint) error
	CreateNotification(notification schema.Notification) (*schema.Notification, error)
	GetUserInfo(c *gin.Context) (types.UserInfoResponse, error)
}

// Ensure PostUtils implements IPostUtils
var _ INotificationUtils = (*NotificationUtils)(nil)

type NotificationUtils struct {
	DB          db.Database
	RedisClient middleware.RedisClientInterface
}

// NewNotificationutils create new notificationUtils
func NewNotificationUtils(DB db.Database, redisClient middleware.RedisClientInterface) *NotificationUtils {
	return &NotificationUtils{
		DB:          DB,
		RedisClient: redisClient,
	}
}

// GetNotificationByUserID
func (nu *NotificationUtils) GetNotificationByUserID(userID uint) ([]schema.Notification, error) {
	var notifications []schema.Notification
	result := nu.DB.Where("user_id = ?", userID).Find(&notifications)
	if result.Error != nil {
		return nil, result.Error
	}
	return notifications, nil
}

// DeleteNotificationByID - delete a notification by id
func (nu *NotificationUtils) DeleteNotificationByID(notificationID uint) error {
	result := nu.DB.Delete(&schema.Notification{}, notificationID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CreateNotification - create a new notification for a user
func (nu *NotificationUtils) CreateNotification(notification schema.Notification) (*schema.Notification, error) {
	result := nu.DB.Create(&notification)
	if result.Error != nil {
		return nil, result.Error
	}
	return &notification, nil
}

// GetUserInfo
func (nu *NotificationUtils) GetUserInfo(c *gin.Context) (types.UserInfoResponse, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL") + "/v1/user/me"

	// Extract the session cookie from the incoming request
	cookie, err := c.Request.Cookie("givegetgo")
	if err != nil {
		return types.UserInfoResponse{}, errors.New("session cookie is missing")
	}

	// Create a new request to the user service
	req, err := http.NewRequest("GET", userServiceURL, nil)
	if err != nil {
		return types.UserInfoResponse{}, err
	}
	req.Header.Set("Cookie", cookie.String()) // Forward the session cookie

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return types.UserInfoResponse{}, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return types.UserInfoResponse{}, errors.New("failed to verify session or session not found")
	}

	// Decode the JSON response into a struct
	var fullResponse types.FullResponseWithData
	if err := json.NewDecoder(resp.Body).Decode(&fullResponse); err != nil {
		return types.UserInfoResponse{}, err
	}

	// Convert the Data field from map to UserInfoResponse
	dataMap, ok := fullResponse.Data.(map[string]interface{})
	if !ok {
		log.Println("Data type assertion to map failed")
		return types.UserInfoResponse{}, fmt.Errorf("response data is not a map")
	}

	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		log.Println("Error marshaling data map to JSON:", err)
		return types.UserInfoResponse{}, err
	}

	var userInfo types.UserInfoResponse
	if err := json.Unmarshal(jsonData, &userInfo); err != nil {
		log.Println("Error unmarshaling JSON to UserInfoResponse:", err)
		return types.UserInfoResponse{}, err
	}

	return userInfo, nil
}
