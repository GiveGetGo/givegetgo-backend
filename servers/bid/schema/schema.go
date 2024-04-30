package schema

import "time"

type BidStatus string

const (
	Submitted BidStatus = "Submitted"
	Accepted  BidStatus = "Accepted"
	Rejected  BidStatus = "Rejected"
)

type PostStatus string

const (
	Active  PostStatus = "Active"
	Matched PostStatus = "Matched"
	Closed  PostStatus = "Closed"
	Expired PostStatus = "Expired"
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

type PostResponse struct {
	PostID      uint       `json:"postID"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Username    string     `json:"username"`
	DatePosted  time.Time  `json:"date_posted"`
	Status      PostStatus `json:"status"`
}
