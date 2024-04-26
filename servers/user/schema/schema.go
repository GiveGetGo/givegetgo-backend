package schema

import (
	"time"
)

type User struct {
	UserID          uint `gorm:"primaryKey"`
	UserName        string
	Email           string
	HashedPassword  string
	Class           string
	Major           string
	ProfileInfo     string
	ReputationScore int
	EmailVerified   bool
	MFAVerified     bool
	MFASecret       string
	DateJoined      time.Time
	LastActiveDate  time.Time
}

type Notification struct {
	NotificationID   uint   `gorm:"primaryKey"`
	UserID           uint   `gorm:"index"`
	NotificationType string `gorm:"type:enum('Bid', 'Match', 'Feedback');default:'Bid'"`
	RelatedID        uint
	Message          string
	DateSent         time.Time
	IsRead           bool
}
