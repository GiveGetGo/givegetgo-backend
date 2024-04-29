package schema

import (
	"time"

	"github.com/GiveGetGo/shared/types"
)

type Notification struct {
	NotificationID   uint `gorm:"primaryKey"`
	UserID           uint `gorm:"index"` // who to send the notification to
	Description      string
	NotificationType types.NotificationType
	CreatedDate      time.Time
}
