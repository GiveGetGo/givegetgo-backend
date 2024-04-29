package schema

import (
	"time"
)

type User struct {
	UserID          uint   `gorm:"primaryKey"`
	UserName        string `gorm:"column:username"`
	Email           string
	HashedPassword  string
	Class           string
	Major           string
	ProfileImage    string
	ProfileInfo     string
	ReputationScore int
	EmailVerified   bool
	MFAVerified     bool
	MFASecret       string
	DateJoined      time.Time
	LastActiveDate  time.Time
}
