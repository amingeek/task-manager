// backend/controllers/progress_additional.go

package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
)

// GetPersonalProgress - دریافت پیشرفت تسک شخصی
func GetPersonalProgress(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var progress models.TaskProgress
	if err := config.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&progress).Error; err != nil {
		// اگر رکوردی وجود ندارد، یک رکورد خالی برگردان
		var task models.Task
		if err := config.DB.First(&task, taskID).Error; err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
			return
		}

		progress = models.TaskProgress{
			TaskID:   task.ID,
			UserID:   userID,
			Progress: 0,
		}
	}
	utils.SuccessResponse(c, http.StatusOK, "OK", progress)
}

// GetGroupProgress - دریافت پیشرفت تسک گروهی برای تمام اعضا
func GetGroupProgress(c *gin.Context) {
	taskID := c.Param("task_id")

	var progress []models.GroupTaskProgress
	if err := config.DB.Where("task_id = ?", taskID).
		Preload("User").
		Find(&progress).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در دریافت پیشرفت")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", progress)
}

// UpdateGroupProgress - بروزرسانی پیشرفت تسک گروهی توسط یک عضو
func UpdateGroupProgress(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("task_id")

	var req UpdateGroupProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// دریافت یا ایجاد رکورد پیشرفت
	var progress models.GroupTaskProgress
	if err := config.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&progress).Error; err != nil {
		// رکورد جدید است
		progress = models.GroupTaskProgress{
			TaskID:     convertStringToUint(taskID),
			UserID:     userID,
			AssignedBy: req.UserID,
			Progress:   req.Progress,
			Notes:      req.Notes,
		}
		config.DB.Create(&progress)
	} else {
		// به‌روزرسانی رکورد موجود
		progress.Progress = req.Progress
		progress.Notes = req.Notes
		config.DB.Save(&progress)
	}

	utils.SuccessResponse(c, http.StatusOK, "پیشرفت با موفقیت بروزرسانی شد", progress)
}
