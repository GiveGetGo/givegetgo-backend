package controller

import (
	"bid/utils"
	"net/http"
	"log"
    "bid/schema"
    "time"
	"strconv"


	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

/*
// func AddBidHandler - add a match
func AddBidHandler(bidUtils *utils.BidUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.LoginRequest
		if err := c.BindJSON(&req); err != nil {
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// TODO: Match
		// for now return success
		types.ResponseSuccess(c, http.StatusCreated, "register", 0, types.UserCreated())
	}
}
*/

type BidRequest struct {
	//thinking UserID may not be neccessary when request a user to post something
	PostID  uint   `json:"postid" binding:"required"`
    UserID  uint   `json:"userid" binding:"required"`
    Description string `json:"description" binding:"required"`
}

func AddBidHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("BidHandler called")

		var req BidRequest
        
        // Binding JSON from the request body
		if err := c.BindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

        // Create a schema.Bid object from the request
        bid := schema.Bid{
            PostID:         uint(req.PostID),
            UserID:         uint(req.UserID),
            BidDescription: req.Description,
            DateSubmitted:  time.Now(),
            Status:         schema.Submitted, // Default status at the time of bid submission
        }

        // Add the bid using the bid utilities
        addedBid, err := bidUtils.AddBid(bid)
        if err != nil {
            log.Printf("Error while adding bid: %v", err)
            types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        log.Printf("Bid added successfully: %v", addedBid)

        // Return the success response with bid creation details
        // Use before pushing PostCreated in my shared document -> change to BidCreated after
        c.JSON(http.StatusOK, gin.H{
            "postid":  req.PostID,
            "Userid":  uint(req.UserID),
            "success": true,
            "message": "Bid added successfully.",
        })
    }
}

// GetBidsForPostHandler handles requests to get all bids for a specific post
func GetBidsForPostHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("GetBidsForPostHandler called")

        postIDParam := c.Param("postID") // The key "PostID" must match the query parameter
        postID, err := strconv.ParseUint(postIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting PostID to uint: %v", err)
            types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        bids, err := bidUtils.GetBidBypostID(uint(postID))
        if err != nil {
            log.Printf("Error fetching bids: %v", err)
            types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        if len(bids) == 0 {
            c.JSON(http.StatusNotFound, gin.H{
                "event": "Bid",
                "code":  "20001",
                "msg":   "No bids found for this post",
            })
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "event": "Bid",
            "postid":    postID,
            "msg":   "Got full bid list for a post",
            "bids":  bids,
        })
    }
}

// FindBidByIDHandler retrieves a specific bid by its ID
func FindBidByIDHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("FindBidByIDHandler called")

        bidIDParam := c.Param("bidID") // The key "bidID" must match the query parameter
        bidID, err := strconv.ParseUint(bidIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting Bidid to uint: %v", err)
            types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        bid, err := bidUtils.GetBidBybidID(uint(bidID))
        if err != nil {
            log.Printf("Error fetching bid: %v", err)
            types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        // Check if the bid was found
        if (len(bid) == 0) {
            c.JSON(http.StatusNotFound, gin.H{
                "event": "Find a bid",
                "code":  20001,
                "id":    bidID,
                "msg":   "No bid found with the specified ID",
            })
            return
        }

        // Return the bid details in the response
        c.JSON(http.StatusOK, gin.H{
            "event": "Find a bid",
            "code":  20000,
            "id":    bidID,
            "msg":   "Found the required bid",
            "data":  bid, // Assuming bid struct fields are correctly tagged for JSON
        })
    }
}

func DeleteBidHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("DeleteBidHandler called")

        // Extract bid ID from query parameter
        bidIDParam := c.Param("bidID")
        bidID, err := strconv.ParseUint(bidIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting bidid to uint: %v", err)
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   true,
                "message": "Invalid request parameter: bidid must be an integer.",
            })
            return
        }

        // Attempt to delete the bid
        err = bidUtils.DeleteBid(uint(bidID))
        if err != nil {
            // Handle the case where the bid could not be deleted
            log.Printf("Error deleting bid: %v", err)
            // You might have different types of errors here, such as "not found" or actual internal errors
            types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        log.Println("Post deleted successfully")
        // Return the success response with deletion confirmation
        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "message": "Bid deleted successfully.",
        })
    }
}

type UpdateBidDescriptionRequest struct {
    Description string `json:"Description"`
}

func UpdateBidDescriptionHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("UpdateBidDescriptionHandler called")

         // Extract bid ID from parameter
         bidIDParam := c.Param("bidID")
         bidID, err := strconv.ParseUint(bidIDParam, 10, 32)
         if err != nil {
             log.Printf("Error converting bidid to uint: %v", err)
             c.JSON(http.StatusBadRequest, gin.H{
                 "error":   true,
                 "message": "Invalid request parameter: bidid must be an integer.",
             })
             return
         }

        // Bind JSON from the request body to UpdateBidDescriptionRequest struct
        var updateReq UpdateBidDescriptionRequest
        if err := c.BindJSON(&updateReq); err != nil {
            log.Printf("Error binding JSON: %v", err)
            types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        // Update the bid description using the bid utilities
        err = bidUtils.UpdateBidDescription(uint(bidID), updateReq.Description)
        if err != nil {
            log.Printf("Error updating bid description: %v", err)
            types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        // Return the success response with bid update details
        c.JSON(http.StatusOK, gin.H{
            "bidid":      bidID,
            "success":    true,
            "message":    "Bid description updated successfully.",
            "newDescription": updateReq.Description,
        })
    }
}