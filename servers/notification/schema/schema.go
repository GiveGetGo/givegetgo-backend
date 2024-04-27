package schema

import (
	"github.com/GiveGetGo/shared/types"
)

type Notification struct {
	NotificationID   uint `gorm:"primaryKey"`
	UserID           uint `gorm:"index"` // who to send the notification to
	NotificationType types.NotificationType
}
