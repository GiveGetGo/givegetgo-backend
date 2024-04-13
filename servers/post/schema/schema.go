package schema

import "time"

type PostStatus string

const (
	Active  PostStatus = "Active"
	Closed  PostStatus = "Closed"
	Expired PostStatus = "Expired"
)

type Post struct {
	PostID      uint `gorm:"primaryKey"`
	UserID      uint `gorm:"index"`
	Title       string
	Description string
	Category    string
	Status      PostStatus
	DatePosted  time.Time
	DateUpdated time.Time
}
