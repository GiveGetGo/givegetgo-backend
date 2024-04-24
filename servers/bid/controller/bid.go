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
	PostID  uint   `json:"postid"`
    UserID  uint   `json:"userid"`
    Message string `json:"message"`
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
            BidDescription: req.Message,
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
        types.ResponseSuccess(c, http.StatusCreated, "bid", addedBid.BidID, types.PostCreated())
    }
}

// GetBidsForPostHandler handles requests to get all bids for a specific post
func GetBidsForPostHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("GetBidsForPostHandler called")

        postIDParam := c.Query("PostID") // The key "PostID" must match the query parameter
        postID, err := strconv.ParseUint(postIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting PostID to uint: %v", err)
            types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        bids, err := bidUtils.GetBidByID(uint(postID))
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

/*
// FindBidByIDHandler retrieves a specific bid by its ID
func FindBidByIDHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("FindBidByIDHandler called")

        bidIDParam := c.Query("Bidid") // The key "PostID" must match the query parameter
        bidID, err := strconv.ParseUint(bidIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting Bidid to uint: %v", err)
            types.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        bid, err := bidUtils.GetBidByID(uint(bidID))
        if err != nil {
            log.Printf("Error fetching bid: %v", err)
            types.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }


        // Check if the bid was found
        if (bid == schema.Bid{}) {
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
*/