package controller

import (
	"errors"
	"net/http"
	"log"
	"time"
	"strconv"
	"post/utils"
	"post/schema"
    "gorm.io/gorm"



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

func AddPostHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("PostHandler called")

		var req types.PostRequest
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
			types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		log.Printf("Post created successfully: %v", addedPost)

		// Return the success response with post creation details
		//Use UserCreated() before pushing my shared document -> change to PostCreated() after
		types.ResponseSuccess(c, http.StatusCreated, "post", uint(req.UserID), types.PostCreated())
	}
}

// GetPostHandler retrieve a specific post by ID using a query parameter.
func GetPostHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("GetPostByIDHandler called")

		// Extract post ID from query parameter
		postIDParam := c.Query("PostID") // The key "PostID" must match the query parameter
		postID, err := strconv.ParseUint(postIDParam, 10, 32)
		if err != nil {
			log.Printf("Error converting PostID to uint: %v", err)
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Retrieve the post using the post utilities
		post, err := postUtils.GetPostByID(uint(postID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error":   true,
					"message": "The requested post was not found.",
				})
			} else {
				log.Printf("Error retrieving post: %v", err)
				types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			}
			return
		}

		log.Printf("Post retrieved successfully: %+v", post)

		// Return the success response with post details
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Post retrieved successfully.",
			"post":    post,
		})
	}
}

func DeletePostHandler(postUtils utils.IPostUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("DeletePostHandler called")

        // Extract post ID from query parameter
        postIDParam := c.Query("postid")
        postID, err := strconv.ParseUint(postIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting postid to uint: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   true,
                "message": "Invalid request parameter: postid must be an integer.",
            })
            return
        }

        // Attempt to delete the post using the post utilities
        err = postUtils.DeletePost(uint(postID))
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                c.JSON(http.StatusNotFound, gin.H{
                    "error": true,
                    "message": "The requested post was not found.",
                })
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": true,
                    "message": "Internal Server Error while deleting the post.",
                })
            }
            return
        }

        log.Println("Post deleted successfully")
        // Return the success response with deletion confirmation
        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "message": "Post deleted successfully.",
        })
    }
}