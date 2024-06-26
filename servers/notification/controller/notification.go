package controller

import (
	"errors"
	"net/http"
	"notification/schema"
	"notification/utils"
	"strconv"
	"time"

	"github.com/GiveGetGo/shared/res"
	"github.com/GiveGetGo/shared/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Public operations
func GetNotification(notificationUtils utils.INotificationUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := notificationUtils.GetUserInfo(c)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		notifications, err := notificationUtils.GetNotificationByUserID(user.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.ResponseError(c, http.StatusNotFound, types.RecordNotFound())
			} else {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			}
			return
		}

		if len(notifications) == 0 {
			res.ResponseError(c, http.StatusNotFound, types.RecordNotFound())
			return
		}

		res.ResponseSuccessWithData(c, http.StatusOK, "Get user notification", types.Success(), notifications)
	}
}

// DeleteNotification - param id is the notification id to be deleted
func DeleteNotification(notificationUtils utils.INotificationUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		notificationIDParam := c.Param("id")
		notificationID, err := strconv.ParseUint(notificationIDParam, 10, 32)
		if err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		err = notificationUtils.DeleteNotificationByID(uint(notificationID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.ResponseError(c, http.StatusNotFound, types.RecordNotFound())
			} else {
				res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			}
			return
		}

		res.ResponseSuccess(c, http.StatusNoContent, "Delete notifications", types.Success())
	}
}

// Internal Operations
func CreateNewNotification(notificationUtils utils.INotificationUtils) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req types.CreateNotificationRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			res.ResponseError(c, http.StatusBadRequest, types.InvalidRequest())
			return
		}

		notification := schema.Notification{
			UserID:           req.UserID,
			Description:      req.Description,
			NotificationType: req.NotificationType,
			CreatedDate:      time.Now(),
		}

		_, err := notificationUtils.CreateNotification(notification)
		if err != nil {
			res.ResponseError(c, http.StatusInternalServerError, types.InternalServerError())
			return
		}

		res.ResponseSuccess(c, http.StatusCreated, "create notification", types.NotificationCreated())
	}
}
