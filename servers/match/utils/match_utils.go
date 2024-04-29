package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"match/db"
	"match/middleware"
	"match/schema"
	"net/http"
	"os"
	"time"

	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

type IMatchUtils interface {
	CreateMatch(postID, postUserID, helperUserID uint) (schema.Match, error)
	GetMatchByID(matchID uint) (schema.Match, error)
	GetAllMatchesByUserID(userid uint) ([]schema.Match, error)
	UpdatePostStatus(postID uint, status schema.PostStatus) error
	DeleteMatch(matchID uint) error
	GetHelperUserID(c *gin.Context, bidId uint) (uint, error)
	GetUserInfo(c *gin.Context) (types.UserInfoResponse, error)
	CreateNotification(userID uint, notificationType types.NotificationType, post schema.PostResponse) error
	GetPostByPostID(c *gin.Context, postID uint) (schema.PostResponse, error)
	FormatNotificationDescription(post schema.PostResponse) string
}

type MatchUtils struct {
	DB          db.Database
	RedisClient middleware.RedisClientInterface
}

// NewMatchUtils creates a new MatchUtils
func NewMatchUtils(DB db.Database, redisClient middleware.RedisClientInterface) *MatchUtils {
	return &MatchUtils{
		DB:          DB,
		RedisClient: redisClient,
	}
}

// func create new match
func (mu *MatchUtils) CreateMatch(postID, postUserID, helperUserID uint) (schema.Match, error) {
	newMatch := schema.Match{
		PostID:       postID,
		PostUserID:   postUserID,
		HelperUserID: helperUserID,
		DateMatched:  time.Now(),
	}

	// Create the match in the database
	result := mu.DB.Create(&newMatch)
	if result.Error != nil {
		return schema.Match{}, result.Error
	}

	// update post status to matched

	return newMatch, nil
}

// func GetMatchByID retrieves a match by its ID
func (mu *MatchUtils) GetMatchByID(matchID uint) (schema.Match, error) {
	var match schema.Match
	err := mu.DB.First(&match, matchID).Error
	if err != nil {
		return schema.Match{}, err
	}

	return match, nil
}

func (mu *MatchUtils) GetAllMatchesByUserID(userid uint) ([]schema.Match, error) {
	var matches []schema.Match

	result := mu.DB.Where("user_id = ?", userid).Find(&matches)
	if result.Error != nil {
		return nil, result.Error
	}

	return matches, nil
}

func (mu *MatchUtils) DeleteMatch(matchID uint) error {
	var match schema.Match
	if result := mu.DB.First(&match, matchID); result.Error != nil {
		return result.Error
	}

	if result := mu.DB.Delete(&match); result.Error != nil {
		return result.Error
	}

	return nil
}

func (mu *MatchUtils) GetHelperUserID(c *gin.Context, bidId uint) (uint, error) {
	bidServiceURL := fmt.Sprintf("%s/v1/bid/%d", os.Getenv("BID_SERVICE_URL"), bidId)

	// Extract the session cookie from the incoming request
	cookie, err := c.Request.Cookie("givegetgo")
	if err != nil {
		return 0, errors.New("session cookie is missing")
	}

	// Create a new request to the user service
	req, err := http.NewRequest("GET", bidServiceURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Cookie", cookie.String()) // Forward the session cookie

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to verify session or session not found")
	}

	// Decode the JSON response into a struct
	var response schema.BidInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, err
	}

	return response.UserID, nil
}

func (mu *MatchUtils) GetUserInfo(c *gin.Context) (types.UserInfoResponse, error) {
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

// UpdatePostStatus sends a request to an external service to update the status of a post
func (mu *MatchUtils) UpdatePostStatus(postID uint, status schema.PostStatus) error {
	// Marshal the request body
	updateReqBody, err := json.Marshal(schema.PostStatusUpdateRequest{
		PostID: postID,
		Status: status,
	})
	if err != nil {
		return err
	}

	// Create the HTTP client and request
	client := &http.Client{}
	postServiceURL := os.Getenv("POST_SERVICE_URL") + "/v1/internal/post/status"
	req, err := http.NewRequest("PUT", postServiceURL, bytes.NewBuffer(updateReqBody))
	if err != nil {
		return err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service", "MATCH")
	req.Header.Set("X-Api-Key", os.Getenv("MATCH_API_KEY"))

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("post service responded with status: %d", resp.StatusCode)
	}

	return nil
}

func (mu *MatchUtils) CreateNotification(userID uint, notificationType types.NotificationType, post schema.PostResponse) error {
	description := mu.FormatNotificationDescription(post)

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
	req.Header.Set("X-Service", "MATCH")
	req.Header.Set("X-Api-Key", os.Getenv("MATCH_API_KEY"))

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

func (mu *MatchUtils) GetPostByPostID(c *gin.Context, postID uint) (schema.PostResponse, error) {
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

func (mu *MatchUtils) FormatNotificationDescription(post schema.PostResponse) string {
	return fmt.Sprintf("Match succeeded with %s for \"%s\". Click in to rate this match!", post.Username, post.Title)
}
