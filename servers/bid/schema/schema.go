package schema

import "time"

type BidStatus string

const (
	Submitted BidStatus = "Submitted"
	Accepted  BidStatus = "Accepted"
	Rejected  BidStatus = "Rejected"
)

type Bid struct {
	BidID          uint `gorm:"primaryKey"`
	PostID         uint `gorm:"index"`
	UserID         uint `gorm:"index"`
	Username       string
	BidDescription string
	DateSubmitted  time.Time
	Status         BidStatus
}

type BidInfoResponse struct {
	UserID         uint   `json:"userID"`
	Username       string `json:"username"`
	BidDescription string `json:"BidDescription"`
	DateSubmitted  string `json:"DateSubmitted"`
}
