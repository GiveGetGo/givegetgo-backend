package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"post/db"
	"post/middleware"
	"post/schema"
	"time"

	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

type IPostUtils interface {
	GetPostByID(postID uint) (schema.Post, error)
	GetPostByUserID(userid uint) ([]schema.Post, error)
	AddPost(post schema.Post) (schema.Post, error)
	GetRecentPosts(days int, count int) ([]schema.Post, error)
	GetArchivePosts(days int, count int) ([]schema.Post, error)
	UpdatePost(postID uint, updateReq types.PostRequest) error
	UpdatePostStatus(postID uint, status schema.PostStatus) error
	DeletePost(postID uint) error
	GetUserInfo(c *gin.Context) (types.UserInfoResponse, error)
}

// Ensure PostUtils implements IPostUtils
var _ IPostUtils = (*PostUtils)(nil)

type PostUtils struct {
	DB          db.Database
	RedisClient middleware.RedisClientInterface
}

// NewPostUtils creates a new PostUtils
func NewPostUtils(DB db.Database, redisClient middleware.RedisClientInterface) *PostUtils {
	return &PostUtils{
		DB:          DB,
		RedisClient: redisClient,
	}
}

// func GetPostByID retrieves a post by its ID
func (pu *PostUtils) GetPostByID(postID uint) (schema.Post, error) {
	var post schema.Post
	err := pu.DB.First(&post, postID).Error
	if err != nil {
		return schema.Post{}, err
	}

	return post, nil
}

func (pu *PostUtils) GetPostByUserID(userid uint) ([]schema.Post, error) {
	var posts []schema.Post
	if err := pu.DB.Where("user_id = ?", userid).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// func AddPost adds a post to the database
func (pu *PostUtils) AddPost(post schema.Post) (schema.Post, error) {
	err := pu.DB.Create(&post).Error
	if err != nil {
		return schema.Post{}, err
	}

	return post, nil
}

func (pu *PostUtils) GetRecentPosts(days int, count int) ([]schema.Post, error) {
	if count == 0 {
		count = 25 // Default value if none is provided and not overridden
	}

	var posts []schema.Post
	since := time.Now().AddDate(0, 0, -days)
	log.Printf("Retrieving posts from the last %d days with a limit of %d posts.", days, count)
	log.Printf("Querying posts since: %v", since)

	result := pu.DB.Where("date_posted >= ? AND status = ?", since, schema.Active).
		Order("date_posted desc").
		Limit(count).
		Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(posts) == 0 {
		log.Println("No posts found with the given criteria")
	} else {
		log.Printf("Found %d posts", len(posts))
	}

	return posts, nil
}

func (pu *PostUtils) GetArchivePosts(days int, count int) ([]schema.Post, error) {
	var posts []schema.Post
	// Calculate the date limit to fetch posts that are older than 'days' days
	since := time.Now().AddDate(0, 0, -days)

	// Adjust the query to exclude posts with 'Active' status
	result := pu.DB.Where("date_posted <= ? AND status != ?", since, schema.Active).
		Order("date_posted desc").
		Limit(count).
		Find(&posts)

	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

func (pu *PostUtils) UpdatePost(postID uint, update types.PostRequest) error {
	// Fetch the existing post from the database
	var post schema.Post
	if err := pu.DB.First(&post, postID).Error; err != nil {
		return err
	}

	// Update fields
	post.PostID = postID
	post.Title = update.Title
	post.Description = update.Description
	post.Category = update.Category
	post.DateUpdated = time.Now()

	return pu.DB.Save(&post).Error
}

// DeletePost deletes a post from the database by its ID.
func (pu *PostUtils) DeletePost(postID uint) error {
	// Attempt to first fetch the post to ensure it exists.
	var post schema.Post
	result := pu.DB.First(&post, postID)
	if result.Error != nil {
		return result.Error
	}

	// If the post exists, proceed to delete it.
	if err := pu.DB.Delete(&post).Error; err != nil {
		return err
	}

	return nil
}

func (pu *PostUtils) GetUserInfo(c *gin.Context) (types.UserInfoResponse, error) {
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

func (pu *PostUtils) UpdatePostStatus(postID uint, status schema.PostStatus) error {
	return pu.DB.Model(&schema.Post{}).Where("post_id = ?", postID).Update("status", status).Error
}
