package schema

import (
	"time"

	"gorm.io/gorm"
)

type RegisterEmailVerification struct {
	gorm.Model
	RegisterVerificationID uint `gorm:"primaryKey"`
	UserID                 uint `gorm:"index"`
	VerificationCode       string
	ExpirationTime         time.Time
	IsVerified             bool `gorm:"default:false"`
}

type ResetPasswordVerification struct {
	gorm.Model
	ResetPasswordID uint `gorm:"primaryKey"`
	UserID          uint `gorm:"index"`
	ResetCode       string
	ExpirationTime  time.Time
	IsUsed          bool `gorm:"default:false"`
}
