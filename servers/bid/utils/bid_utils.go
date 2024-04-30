package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"bid/db"
	"bid/middleware"
	"bid/schema"

	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

type IBidUtils interface {
	GetBidBypostID(bidID uint) ([]schema.Bid, error)
	AddBid(bid schema.Bid) (schema.Bid, error)
	GetBidBybidID(bidID uint) ([]schema.Bid, error)
	DeleteBid(bidID uint) error
	UpdateBidDescription(bidID uint, description string) error
	GetUserInfo(c *gin.Context) (types.UserInfoResponse, error)
	CreateNotification(userID uint, notificationType types.NotificationType, post schema.PostResponse) error
	GetPostByPostID(c *gin.Context, postID uint) (schema.PostResponse, error)
	FormatNotificationDescription(post schema.PostResponse) string
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
		return result.Error // Return the error (e.g., not found)
	}

	// If the post exists, proceed to delete it.
	if err := pu.DB.Delete(&bid).Error; err != nil {
		return err // Return any error that occurs during the delete operation.
	}

	return nil // Return nil if the delete operation is successful.
}

// UpdateBidDescription updates the description of a bid identified by bidID
func (bu *BidUtils) UpdateBidDescription(bidID uint, description string) error {
	// Find the bid by ID
	var bid schema.Bid
	result := bu.DB.First(&bid, "bid_id = ?", bidID)
	if result.Error != nil {
		return result.Error // If not found or other DB error
	}

	// Update the bid's description
	bid.BidDescription = description
	return bu.DB.Save(&bid).Error
}

func (bu *BidUtils) GetUserInfo(c *gin.Context) (types.UserInfoResponse, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL") + "/v1/user/me"

	// Extract the session cookie from the incoming request
	cookie, err := c.Request.Cookie("givegetgo")
	if err != nil {
		return types.UserInfoResponse{}, errors.New("session cookie is missing")
	}

	// Create a new request to the user service
	req, err := http.NewRequest("GET", userServiceURL, nil)
	if err != nil {
		return types.UserInfoResponse{}, err
	}
	req.Header.Set("Cookie", cookie.String()) // Forward the session cookie

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return types.UserInfoResponse{}, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return types.UserInfoResponse{}, errors.New("failed to verify session or session not found")
	}

	// Decode the JSON response into a struct
	var fullResponse types.FullResponseWithData
	if err := json.NewDecoder(resp.Body).Decode(&fullResponse); err != nil {
		return types.UserInfoResponse{}, err
	}

	// Convert the Data field from map to UserInfoResponse
	dataMap, ok := fullResponse.Data.(map[string]interface{})
	if !ok {
		log.Println("Data type assertion to map failed")
		return types.UserInfoResponse{}, fmt.Errorf("response data is not a map")
	}

	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		log.Println("Error marshaling data map to JSON:", err)
		return types.UserInfoResponse{}, err
	}

	var userInfo types.UserInfoResponse
	if err := json.Unmarshal(jsonData, &userInfo); err != nil {
		log.Println("Error unmarshaling JSON to UserInfoResponse:", err)
		return types.UserInfoResponse{}, err
	}

	return userInfo, nil
}

func (bu *BidUtils) CreateNotification(userID uint, notificationType types.NotificationType, post schema.PostResponse) error {
	description := bu.FormatNotificationDescription(post)

	// Marshal the request body
	notificationReqBody, err := json.Marshal(types.CreateNotificationRequest{
		UserID:           userID,
		Description:      description,
		NotificationType: notificationType,
	})
	if err != nil {
		return err
	}

	// Create the HTTP client and request
	client := &http.Client{}
	notificationServiceURL := os.Getenv("NOTIFICATION_SERVICE_URL") + "/v1/internal/notification"
	req, err := http.NewRequest("POST", notificationServiceURL, bytes.NewBuffer(notificationReqBody))
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service", "BID")
	req.Header.Set("X-Api-Key", os.Getenv("BID_API_KEY"))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("notification service responded with status: %d", resp.StatusCode)
	}

	return nil
}

func (bu *BidUtils) GetPostByPostID(c *gin.Context, postID uint) (schema.PostResponse, error) {
	postServiceURL := os.Getenv("POST_SERVICE_URL") + fmt.Sprintf("/v1/post/%d", postID)

	// Extract the session cookie from the incoming request
	cookie, err := c.Request.Cookie("givegetgo")
	if err != nil {
		return schema.PostResponse{}, errors.New("session cookie is missing")
	}

	// Create a new request to the post service
	req, err := http.NewRequest("GET", postServiceURL, nil)
	if err != nil {
		return schema.PostResponse{}, err
	}
	req.Header.Set("Cookie", cookie.String()) // Forward the session cookie

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return schema.PostResponse{}, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return schema.PostResponse{}, fmt.Errorf("failed to retrieve post or post not found, status code: %d", resp.StatusCode)
	}

	// Decode the JSON response into a struct
	var fullResponse types.FullResponseWithData
	if err := json.NewDecoder(resp.Body).Decode(&fullResponse); err != nil {
		return schema.PostResponse{}, err
	}

	// Convert the Data field from map to PostResponse
	dataMap, ok := fullResponse.Data.(map[string]interface{})
	if !ok {
		log.Println("Data type assertion to map failed")
		return schema.PostResponse{}, fmt.Errorf("response data is not a map")
	}

	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		log.Println("Error marshaling data map to JSON:", err)
		return schema.PostResponse{}, err
	}

	var postResponse schema.PostResponse
	if err := json.Unmarshal(jsonData, &postResponse); err != nil {
		log.Println("Error unmarshaling JSON to PostResponse:", err)
		return schema.PostResponse{}, err
	}

	return postResponse, nil
}

func (bu *BidUtils) FormatNotificationDescription(post schema.PostResponse) string {
	return fmt.Sprintf("New matching request %s for \"%s\".", post.Username, post.Title)
}
