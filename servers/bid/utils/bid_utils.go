package utils

import (
	"bid/db"
	"bid/middleware"
	"bid/schema"
)

type IBidUtils interface {
	GetBidByID(bidID uint) ([]schema.Bid, error)
	AddBid(bid schema.Bid) (schema.Bid, error)
}

// Ensure PostUtils implements IPostUtils
var _ IBidUtils = (*BidUtils)(nil)

type BidUtils struct {
	DB          db.Database
	RedisClient middleware.RedisClientInterface
}

// NewbidUtils creates a new bidUtils
func NewBidUtils(DB db.Database, redisClient middleware.RedisClientInterface) *BidUtils {
	return &BidUtils{
		DB:          DB,
		RedisClient: redisClient,
	}
}

// func GetBidByID retrieves a bid by its ID
func (bu *BidUtils) GetBidByID(postID uint) ([]schema.Bid, error) {
    var bids []schema.Bid
    err := bu.DB.Where("post_id = ?", postID).Find(&bids).Error
    if err != nil {
        return nil, err
    }
    return bids, nil
}

// func Addbid adds a bid to the database
func (bu *BidUtils) AddBid(bid schema.Bid) (schema.Bid, error) {
	err := bu.DB.Create(&bid).Error
	if err != nil {
		return schema.Bid{}, err
	}

	return bid, nil
}
