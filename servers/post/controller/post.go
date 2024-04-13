package controller

import (
	"net/http"
	"log"
	"time"
	"post/utils"
	"post/schema"

	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

/*
// AddPostHandler is a function that handles the request to add a post
func AddPostHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.LoginRequest
		if err := c.BindJSON(&req); err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// TODO: Post validation
		// for now return success
		types.ResponseSuccess(c, http.StatusCreated, "register", 0, types.UserCreated())
	}
}
*/

type PostRequest struct {
	UserID      int64  `json:"userID" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

// Response structure as per the API specification
type ApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func AddPostHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("PostHandler called")

		var req PostRequest
		if err := c.BindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Create a schema.Post object from the request
		post := schema.Post{
			UserID:      uint(req.UserID),
			Title:       req.Title,
			Description: req.Description,
			Category:    req.Category,
			Status:      schema.Active, 
			DatePosted:  time.Now(),
			DateUpdated: time.Now(),
		}


		// Add the post using the post utilities
		addedPost, err := postUtils.AddPost(post)
		if err != nil {
			log.Printf("Error creating post: %v", err)
			c.JSON(http.StatusInternalServerError, ApiResponse{
				Code: 50001,
				Msg:  "Internal Server Error",
			})
			return
		}

		log.Printf("Post created successfully: %v", addedPost)

		// Return the success response with post creation details
		//Use UserCreated() before pushing my shared document -> change to PostCreated() after
		types.ResponseSuccess(c, http.StatusCreated, "post", 0, types.UserCreated())
	}
}
