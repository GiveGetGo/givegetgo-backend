package utils

import (
	"notification/db"
	"notification/middleware"
	"notification/schema"

	"gorm.io/gorm"
)

type INotificationUtils interface {
	GetNotificationByUserID(userID uint) ([]schema.Notification, error)
	DeleteNotificationByID(notificationID uint) error
	CreateNotification(notification schema.Notification) (*schema.Notification, error)
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
