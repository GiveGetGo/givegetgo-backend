package utils

import (
	"post/db"
	"post/middleware"
	"post/schema"
)

type IPostUtils interface {
	GetPostByID(postID uint) (schema.Post, error)
	AddPost(post schema.Post) (schema.Post, error)
	DeletePost(postID uint) error

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

// func AddPost adds a post to the database
func (pu *PostUtils) AddPost(post schema.Post) (schema.Post, error) {
	err := pu.DB.Create(&post).Error
	if err != nil {
		return schema.Post{}, err
	}

	return post, nil
}

// DeletePost deletes a post from the database by its ID.
func (pu *PostUtils) DeletePost(postID uint) error {
    // Attempt to first fetch the post to ensure it exists.
    var post schema.Post
    result := pu.DB.First(&post, postID)
    if result.Error != nil {
        return result.Error  // Return the error (e.g., not found)
    }

    // If the post exists, proceed to delete it.
    if err := pu.DB.Delete(&post).Error; err != nil {
        return err  // Return any error that occurs during the delete operation.
    }

    return nil  // Return nil if the delete operation is successful.
}