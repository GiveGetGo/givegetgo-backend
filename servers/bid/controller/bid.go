package controller

import (
	"bid/utils"
	"net/http"
	"log"
    "bid/schema"
    "time"

	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

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

/*
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
        types.ResponseSuccess(c, http.StatusCreated, "bid", addedBid.BidID, types.BidCreated())
    }
}
*/