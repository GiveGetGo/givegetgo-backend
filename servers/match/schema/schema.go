package schema

import "time"

type MatchStatus string

const (
	Matching    MatchStatus = "Matching"
	Matched     MatchStatus = "Matched"
	Unfulfilled MatchStatus = "Unfulfilled"
)

type Match struct {
	MatchID            uint `gorm:"primaryKey"`
	PostID             uint `gorm:"index"`
	HelperUserID       uint `gorm:"index"`
	Status             MatchStatus
	DateMatched        time.Time
	FulfillmentDetails string
}
