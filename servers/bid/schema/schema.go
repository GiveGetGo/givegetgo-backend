package schema

import "time"

type BidStatus string

const (
	Submitted BidStatus = "Submitted"
	Accepted  BidStatus = "Accepted"
	Rejected  BidStatus = "Rejected"
)

type Bid struct {
	BidID           uint `gorm:"primaryKey"`
	PostID          uint `gorm:"index"`
	UserID          uint `gorm:"index"`
	BidDescription  string
	TermsConditions string
	DateSubmitted   time.Time
	Status          BidStatus
}
