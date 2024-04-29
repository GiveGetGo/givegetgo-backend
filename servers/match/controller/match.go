package controller

import (
	"match/schema"
	"match/utils"
	"net/http"
	"strconv"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
)

// func MatchHandler - add a match
func MatchHandler(matchUtils utils.IMatchUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.MatchRequest
		if err := c.BindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		user, err := matchUtils.GetUserInfo(c)
		if err != nil {
			res.ResponseError(c, http.StatusUnauthorized, types.InvalidCredentials())
			return
		}

		helperUserid, err := matchUtils.GetHelperUserID(c, req.BidID)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		post, err := matchUtils.GetPostByPostID(c, req.PostID)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// create new match
		_, err = matchUtils.CreateMatch(post.PostID, user.UserID, helperUserid)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// after created match, update post status to matched
		err = matchUtils.UpdatePostStatus(req.PostID, schema.Matched)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Notify both post user and bid user
		err = matchUtils.CreateNotification(user.UserID, types.NewMatch, post)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		err = matchUtils.CreateNotification(helperUserid, types.BidMatch, post)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "new-match", types.MatchSuccess())
	}
}

// GetMatchHandler - get match by matchid
func GetMatchHandler(matchUtils utils.IMatchUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract match ID from URL
		matchIDParam := c.Param("id")
		matchID, err := strconv.ParseUint(matchIDParam, 10, 32)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Retrieve the match by ID
		match, err := matchUtils.GetMatchByID(uint(matchID))
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		// Return the match data
		res.ResponseSuccessWithData(c, http.StatusOK, "get-match", types.Success(), match)
	}
}

// DeleteMatchHandler - delete match by matchid
func DeleteMatchHandler(matchUtils utils.IMatchUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract match ID from URL
		matchIDParam := c.Param("id")
		matchID, err := strconv.ParseUint(matchIDParam, 10, 32)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		// Delete the match using the match utility
		if err := matchUtils.DeleteMatch(uint(matchID)); err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusOK, "delete match", types.Success())
	}
}
