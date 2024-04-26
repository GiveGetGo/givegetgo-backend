package controller

import (
	"errors"
	"log"
	"net/http"
	"post/schema"
	"post/utils"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

func AddPostHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.PostRequest
		if err := c.BindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		user, err := postUtils.GetUserInfo(c)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Create a schema.Post object from the request
		post := schema.Post{
			UserID:      user.UserID,
			Username:    user.Username,
			Title:       req.Title,
			Description: req.Description,
			Category:    req.Category,
			Status:      schema.Active,
			DatePosted:  time.Now(),
			DateUpdated: time.Now(),
		}

		// Add the post using the post utilities
		_, err = postUtils.AddPost(post)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return the success response with post creation details
		//Use UserCreated() before pushing my shared document -> change to PostCreated() after
		res.ResponseSuccess(c, http.StatusCreated, "post", types.PostCreated())
	}
}

// GetPostHandler retrieve post from recent
func GetPostHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set default values
		days := 14 // two weeks
		var count int

		// Check if day is specified in the query
		if dayParam, ok := c.GetQuery("day"); ok {
			day, err := strconv.Atoi(dayParam)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}
			if day < 1 {
				days = 1
			} else if day > 30 {
				days = 30
			} else {
				days = day
			}
		}

		// Check if count is specified in the query
		if countParam, ok := c.GetQuery("count"); ok {
			cnt, err := strconv.Atoi(countParam)
			if err != nil {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
				return
			}
			if cnt > 0 {
				count = cnt
			}
		}

		// Retrieve posts from the recent days
		posts, err := postUtils.GetRecentPosts(days, count)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Convert posts to the response structure
		var responsePosts []schema.PostResponse
		for _, post := range posts {
			responsePosts = append(responsePosts, schema.PostResponse{
				PostID:      post.PostID,
				Title:       post.Title,
				Description: post.Description,
				Category:    post.Category,
				Username:    post.Username,
				DatePosted:  post.DatePosted,
				Status:      post.Status,
			})
		}

		// Return the success response with post details
		res.ResponseSuccessWithData(c, http.StatusOK, "Post retrieved", types.Success(), responsePosts)
	}
}

func GetPostArchiveHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set default values
		days := 28 // four weeks
		var count int

		// Check if day is specified in the query
		if dayParam, ok := c.GetQuery("day"); ok {
			day, err := strconv.Atoi(dayParam)
			if err != nil {
				res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
				return
			}
			if day < 1 {
				days = 1
			} else {
				days = day
			}
		}

		// Check if count is specified in the query
		if countParam, ok := c.GetQuery("count"); ok {
			cnt, err := strconv.Atoi(countParam)
			if err != nil {
				res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
				return
			}
			if cnt > 0 {
				count = cnt
			}
		}

		// Retrieve posts from the recent days
		posts, err := postUtils.GetArchivePosts(days, count)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Convert posts to the response structure
		var responsePosts []schema.PostResponse
		for _, post := range posts {
			responsePosts = append(responsePosts, schema.PostResponse{
				PostID:      post.PostID,
				Title:       post.Title,
				Description: post.Description,
				Category:    post.Category,
				Username:    post.Username,
				DatePosted:  post.DatePosted,
				Status:      post.Status,
			})
		}

		// Return the success response with post details
		res.ResponseSuccessWithData(c, http.StatusOK, "Post retrieved", types.Success(), responsePosts)
	}
}

func GetPostByPostIdHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDParam := c.Param("id")
		postID, err := strconv.ParseUint(postIDParam, 10, 32)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		post, err := postUtils.GetPostByID(uint(postID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.ResponseError(c, http.StatusNotFound, types.RecordNotFound())
			} else {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			}
			return
		}

		responsePost := schema.PostResponse{
			PostID:      post.PostID,
			Title:       post.Title,
			Description: post.Description,
			Category:    post.Category,
			Username:    post.Username,
			DatePosted:  post.DatePosted,
			Status:      post.Status,
		}

		res.ResponseSuccessWithData(c, http.StatusOK, "Post Retrieved", types.Success(), responsePost)
	}
}

func GetPostByUserIdHandler(postutils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get user id from user me endpoint
		user, err := postutils.GetUserInfo(c)
		if err != nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		posts, err := postutils.GetPostByUserID(user.UserID)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Convert posts to the response structure
		var responsePosts []schema.PostResponse
		for _, post := range posts {
			responsePosts = append(responsePosts, schema.PostResponse{
				PostID:      post.PostID,
				Title:       post.Title,
				Description: post.Description,
				Category:    post.Category,
				Username:    post.Username,
				DatePosted:  post.DatePosted,
				Status:      post.Status,
			})
		}

		res.ResponseSuccessWithData(c, http.StatusOK, "get post by user", types.Success(), responsePosts)
	}
}

func EditPostByIdHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDParam := c.Param("id")
		postID, err := strconv.ParseUint(postIDParam, 10, 32)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		var updateReq types.PostRequest
		if err := c.BindJSON(&updateReq); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		err = postUtils.UpdatePost(uint(postID), updateReq)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "edit post", types.Success())
	}
}

func DeletePostHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract post ID from query parameter
		postIDParam := c.Param("id")
		postID, err := strconv.ParseUint(postIDParam, 10, 32)
		if err != nil {
			log.Printf("Error converting postid to uint: %v", err)
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Attempt to delete the post using the post utilities
		err = postUtils.DeletePost(uint(postID))
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return the success response with deletion confirmation
		res.ResponseSuccess(c, http.StatusOK, "delete post", types.Success())
	}
}

func UpdatePostStatusHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the request body
		var updateReq schema.PostStatusUpdateRequest
		if err := c.BindJSON(&updateReq); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Validate the new status
		if updateReq.Status != schema.Active && updateReq.Status != schema.Matched &&
			updateReq.Status != schema.Closed && updateReq.Status != schema.Expired {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Update the post status
		if err := postUtils.UpdatePostStatus(uint(updateReq.PostID), updateReq.Status); err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "update post sucess", types.Success())
	}
}
