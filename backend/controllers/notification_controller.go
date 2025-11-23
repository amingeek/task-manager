// backend/controllers/notification_controller.go

package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
)

// GetNotifications - دریافت تمام اعلان‌های کاربر
func GetNotifications(c *gin.Context) {
	userID := c.GetUint("userID")

	var notifications []models.Notification
	if err := config.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در دریافت اعلان‌ها")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", notifications)
}

// MarkAsRead - علامت‌گذاری اعلان به عنوان خوانده شده
func MarkAsRead(c *gin.Context) {
	userID := c.GetUint("userID")
	notificationID := c.Param("id")

	// بررسی اینکه اعلان متعلق به کاربر است
	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "اعلان پیدا نشد")
		return
	}

	// علامت‌گذاری به عنوان خوانده شده
	notification.IsRead = true
	if err := config.DB.Save(&notification).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در به‌روزرسانی اعلان")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "اعلان به عنوان خوانده شده علامت‌گذاری شد", notification)
}

// DeleteNotification - حذف اعلان
func DeleteNotification(c *gin.Context) {
	userID := c.GetUint("userID")
	notificationID := c.Param("id")

	// بررسی اینکه اعلان متعلق به کاربر است
	var notification models.Notification
	if err := config.DB.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "اعلان پیدا نشد")
		return
	}

	// حذف اعلان
	if err := config.DB.Delete(&notification).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در حذف اعلان")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "اعلان با موفقیت حذف شد", nil)
}
