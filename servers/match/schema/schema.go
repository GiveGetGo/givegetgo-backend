package schema

import "time"

// type MatchStatus string
// const (
// 	Matching    MatchStatus = "Matching"
// 	Matched     MatchStatus = "Matched"
// 	Unfulfilled MatchStatus = "Unfulfilled"
// )

type PostStatus string

const (
	Active  PostStatus = "Active"
	Matched PostStatus = "Matched"
	Closed  PostStatus = "Closed"
	Expired PostStatus = "Expired"
)

type Match struct {
	MatchID        uint `gorm:"primaryKey"`
	PostID         uint `gorm:"index"`
	PostUserID     uint `gorm:"index"`
	HelperUserID   uint `gorm:"index"`
	PostUsername   string
	HelperUsername string
	// Status             MatchStatus
	DateMatched        time.Time
	FulfillmentDetails string
}

type PostStatusUpdateRequest struct {
	PostID uint       `json:"postID"`
	Status PostStatus `json:"status"`
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
