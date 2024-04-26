package utils

import (
	"bid/db"
	"bid/middleware"
	"bid/schema"
)

type IBidUtils interface {
	GetBidBypostID(bidID uint) ([]schema.Bid, error)
	AddBid(bid schema.Bid) (schema.Bid, error)
	GetBidBybidID(bidID uint) ([]schema.Bid, error)
	DeleteBid(bidID uint) error
    UpdateBidDescription(bidID uint, description string) error
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

// func GetBidBypostID retrieves a bid by its postID
func (bu *BidUtils) GetBidBypostID(postID uint) ([]schema.Bid, error) {
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

// func GetBidBypostID retrieves a bid by its bidID
func (bu *BidUtils) GetBidBybidID(bidID uint) ([]schema.Bid, error) {
    var bids []schema.Bid
    err := bu.DB.Where("bid_id = ?", bidID).Find(&bids).Error
    if err != nil {
        return nil, err
    }
    return bids, nil
}

// DeletePost deletes a post from the database by its ID.
func (pu *BidUtils) DeleteBid(bidID uint) error {
    // Attempt to first fetch the post to ensure it exists.
    var bid schema.Bid
    result := pu.DB.First(&bid, bidID)
    if result.Error != nil {
        return result.Error  // Return the error (e.g., not found)
    }

    // If the post exists, proceed to delete it.
    if err := pu.DB.Delete(&bid).Error; err != nil {
        return err  // Return any error that occurs during the delete operation.
    }

    return nil  // Return nil if the delete operation is successful.
}

// UpdateBidDescription updates the description of a bid identified by bidID
func (bu *BidUtils) UpdateBidDescription(bidID uint, description string) error {
    // Find the bid by ID
    var bid schema.Bid
    result := bu.DB.First(&bid, "bid_id = ?", bidID)
    if result.Error != nil {
        return result.Error  // If not found or other DB error
    }

    // Update the bid's description
    bid.BidDescription = description
    return bu.DB.Save(&bid).Error
}