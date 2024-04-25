package schema

import "time"

type PostStatus string

const (
	Active  PostStatus = "Active"
	Matched PostStatus = "Matched"
	Closed  PostStatus = "Closed"
	Expired PostStatus = "Expired"
)

type Post struct {
	PostID      uint `gorm:"primaryKey"`
	UserID      uint `gorm:"index"`
	Username    string
	Title       string
	Description string
	Category    string
	Status      PostStatus
	DatePosted  time.Time
	DateUpdated time.Time
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

type PostStatusUpdateRequest struct {
	PostID uint       `json:"postID"`
	Status PostStatus `json:"status"`
}
