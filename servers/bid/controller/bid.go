package controller

import (
	"bid/utils"
	"net/http"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

// func AddBidHandler - add a match
func AddBidHandler(bidUtils *utils.BidUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.LoginRequest
		if err := c.BindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// TODO: Match
		// for now return success
		res.ResponseSuccess(c, http.StatusCreated, "register", types.UserCreated())
	}
}
