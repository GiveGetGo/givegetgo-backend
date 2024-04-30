package controller

import (
	"bid/utils"
	"net/http"
	"log"
    "bid/schema"
    "time"
	"strconv"

	"github.com/GiveGetGo/shared/res"
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

func AddBidHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("BidHandler called")

		var req types.BidRequest
        
        // Binding JSON from the request body
        if err := c.BindJSON(&req); err != nil {
			log.Printf("Error binding JSON: %v", err)
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

        postIDParam := c.Param("postid") // The key "PostID" must match the query parameter
        postID, err := strconv.ParseUint(postIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting PostID to uint: %v", err)
            res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        user, err := bidUtils.GetUserInfo(c)
		if err != nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
        }

        // Create a schema.Bid object from the request
        bid := schema.Bid{
            PostID:         uint(postID),
            UserID:         user.UserID,
            Username:       user.Username,
            BidDescription: req.Description,
            DateSubmitted:  time.Now(),
            Status:         schema.Submitted, // Default status at the time of bid submission
        }

        // Add the bid using the bid utilities
        addedBid, err := bidUtils.AddBid(bid)
        if err != nil {
            log.Printf("Error while adding bid: %v", err)
            res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        post, err := bidUtils.GetPostByPostID(c, uint(postID))
        if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

        // Notify the bid user
		err = bidUtils.CreateNotification(user.UserID, types.NewBid, post)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

        log.Printf("Bid added successfully: %v", addedBid)

        // Return the success response with bid creation details
        // Use before pushing PostCreated in my shared document -> change to BidCreated after

        res.ResponseSuccess(c, http.StatusCreated, "bid", types.BidCreated())
    }
}

// GetBidsForPostHandler handles requests to get all bids for a specific post
func GetBidsForPostHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("GetBidsForPostHandler called")

        // get user id from user me endpoint
        user, err := bidUtils.GetUserInfo(c)
		if err != nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
        }

        postIDParam := c.Param("postid") // The key "PostID" must match the query parameter
        postID, err := strconv.ParseUint(postIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting PostID to uint: %v", err)
            res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        bids, err := bidUtils.GetBidBypostID(uint(postID))
        if err != nil {
            log.Printf("Error fetching bids: %v", err)
            res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        //if there is no bid with the postID
        if len(bids) == 0 {
            res.ResponseError(c, http.StatusUnauthorized, types.InvalidRequest())
            return
        }

        // Convert bids to the response structure
		var responseBids []schema.BidInfoResponse
		for _, bid := range bids {
			responseBids = append(responseBids, schema.BidInfoResponse{
				UserID:         user.UserID,
                Username:       user.Username,
                BidDescription: bid.BidDescription,
                DateSubmitted:  bid.DateSubmitted.Format(time.RFC3339),
			})
		}

		res.ResponseSuccessWithData(c, http.StatusOK, "Got full bid list for a post", types.Success(), responseBids)
    }
}

// FindBidByIDHandler retrieves a specific bid by its ID
func FindBidByIDHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("FindBidByIDHandler called")

        // get user id from user me endpoint
        user, err := bidUtils.GetUserInfo(c)
		if err != nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
        }

        bidIDParam := c.Param("bidid") // The key "bidID" must match the query parameter
        bidID, err := strconv.ParseUint(bidIDParam, 10, 32)
        if err != nil {
            log.Printf("Error converting Bidid to uint: %v", err)
            res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        bid, err := bidUtils.GetBidBybidID(uint(bidID))
        if err != nil {
            log.Printf("Error fetching bid: %v", err)
            res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        // Check if the bid was found
        if (len(bid) == 0) {
            res.ResponseError(c, http.StatusUnauthorized, types.InvalidRequest())
            return
        }

        // Convert bids to the response structure
		var responseBids []schema.BidInfoResponse
		for _, bid := range bid {
			responseBids = append(responseBids, schema.BidInfoResponse{
				UserID:         user.UserID,
                Username:       user.Username,
                BidDescription: bid.BidDescription,
                DateSubmitted:  bid.DateSubmitted.Format(time.RFC3339),
			})
		}

		res.ResponseSuccessWithData(c, http.StatusOK, "Find a bid", types.Success(), responseBids)
    
    }
}

func DeleteBidHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("DeleteBidHandler called")

        // Extract bid ID from query parameter
        bidIDParam := c.Param("bidid")
        bidID, err := strconv.ParseUint(bidIDParam, 10, 32)
        if err != nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
        }

        // Attempt to delete the bid
        err = bidUtils.DeleteBid(uint(bidID))
        if err != nil {
            // Handle the case where the bid could not be deleted
            log.Printf("Error deleting bid: %v", err)
            // You might have different types of errors here, such as "not found" or actual internal errors
            res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        log.Println("Post deleted successfully")
        // Return the success response with deletion confirmation
        res.ResponseSuccess(c, http.StatusOK, "delete bid", types.Success())
    }
}

type UpdateBidDescriptionRequest struct {
    Description string `json:"Description"`
}

func UpdateBidDescriptionHandler(bidUtils utils.IBidUtils) gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("UpdateBidDescriptionHandler called")

         // Extract bid ID from parameter
         bidIDParam := c.Param("bidid")
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
            res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
            return
        }

        // Update the bid description using the bid utilities
        err = bidUtils.UpdateBidDescription(uint(bidID), updateReq.Description)
        if err != nil {
            log.Printf("Error updating bid description: %v", err)
            res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
            return
        }

        res.ResponseSuccessWithData(c, http.StatusOK, "update bid", types.Success(), updateReq.Description)

    }
}