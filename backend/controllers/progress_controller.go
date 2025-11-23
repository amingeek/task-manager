// backend/controllers/progress_controller.go

package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateProgressRequest struct {
	Progress int    `json:"progress" binding:"required,min=0,max=100"`
	Notes    string `json:"notes"`
}

type UpdateGroupProgressRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Progress int    `json:"progress" binding:"required,min=0,max=100"`
	Notes    string `json:"notes"`
	Approved bool   `json:"approved"`
}

// UpdatePersonalProgress - بروزرسانی پیشرفت تسک شخصی
func UpdatePersonalProgress(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// بررسی وجود تسک و مالکیت
	var task models.Task
	if err := config.DB.Where("id = ? AND creator_id = ? AND is_group_task = ?", taskID, userID, false).First(&task).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "تسک شخصی پیدا نشد")
		return
	}

	// یافتن یا ایجاد رکورد پیشرفت
	var progress models.TaskProgress
	if err := config.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&progress).Error; err != nil {
		progress = models.TaskProgress{
			TaskID:   task.ID,
			UserID:   userID,
			Progress: req.Progress,
			Notes:    req.Notes,
		}
		config.DB.Create(&progress)
	} else {
		progress.Progress = req.Progress
		progress.Notes = req.Notes
		config.DB.Save(&progress)
	}

	// بررسی تکمیل شدن تسک
	if req.Progress == 100 {
		progress.IsCompleted = true
		now := time.Now()
		progress.CompletedAt = &now
		// آپدیت وضعیت تسک
		task.Status = models.StatusCompleted
		config.DB.Save(&task)
	} else if req.Progress > 0 {
		task.Status = models.StatusInProgress
		config.DB.Save(&task)
	}

	utils.SuccessResponse(c, http.StatusOK, "پیشرفت با موفقیت بروزرسانی شد", progress)
}

// UpdateGroupMemberProgress - بروزرسانی پیشرفت عضو توسط مدیر گروه
func UpdateGroupMemberProgress(c *gin.Context) {
	userID := c.GetUint("userID") // کاربر مدیر
	taskID := c.Param("id")

	var req UpdateGroupProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// بررسی تسک و دسترسی مدیر
	var task models.Task
	if err := config.DB.Preload("Group").First(&task, taskID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
		return
	}

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", task.GroupID, userID, "admin").First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه میتوانند پیشرفت را بروزرسانی کنند")
		return
	}

	// یافتن یا ایجاد رکورد پیشرفت
	var progress models.GroupTaskProgress
	if err := config.DB.Where("task_id = ? AND user_id = ?", taskID, req.UserID).First(&progress).Error; err != nil {
		progress = models.GroupTaskProgress{
			TaskID:     task.ID,
			UserID:     req.UserID,
			AssignedBy: userID,
			Progress:   req.Progress,
			Notes:      req.Notes,
			Approved:   req.Approved,
		}
		config.DB.Create(&progress)
	} else {
		progress.Progress = req.Progress
		progress.Notes = req.Notes
		progress.AssignedBy = userID
		if req.Approved && !progress.Approved {
			now := time.Now()
			progress.Approved = true
			progress.ApprovedBy = &userID
			progress.ApprovedAt = &now
		}
		config.DB.Save(&progress)
	}

	// بررسی تکمیل شدن
	if req.Progress == 100 {
		progress.IsCompleted = true
		now := time.Now()
		progress.CompletedAt = &now
		config.DB.Save(&progress)
	}

	// بررسی اینکه آیا همه اعضا تکمیل کردهاند
	var allProgress []models.GroupTaskProgress
	config.DB.Where("task_id = ?", taskID).Find(&allProgress)

	allCompleted := true
	for _, p := range allProgress {
		if p.Progress < 100 {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		task.Status = models.StatusCompleted
		config.DB.Save(&task)
	}

	utils.SuccessResponse(c, http.StatusOK, "پیشرفت عضو با موفقیت بروزرسانی شد", progress)
}

// GetTaskProgress - دریافت پیشرفت تسک (شخصی یا گروهی)
func GetTaskProgress(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
		return
	}

	if task.IsGroupTask {
		// برای تسک گروهی
		var progress []models.GroupTaskProgress
		config.DB.Where("task_id = ?", taskID).Preload("User").Find(&progress)
		utils.SuccessResponse(c, http.StatusOK, "OK", progress)
	} else {
		// برای تسک شخصی
		var progress models.TaskProgress
		if err := config.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&progress).Error; err != nil {
			// ایجاد رکورد خالی اگر وجود ندارد
			progress = models.TaskProgress{
				TaskID:   task.ID,
				UserID:   userID,
				Progress: 0,
			}
		}
		utils.SuccessResponse(c, http.StatusOK, "OK", progress)
	}
}

// GetMyGroupProgress - دریافت پیشرفت شخصی در تسک گروهی
func GetMyGroupProgress(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var progress models.GroupTaskProgress
	if err := config.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&progress).Error; err != nil {
		// اگر رکوردی وجود ندارد، یک رکورد خالی برگردان
		var task models.Task
		if err := config.DB.First(&task, taskID).Error; err != nil {
			utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
			return
		}

		progress = models.GroupTaskProgress{
			TaskID:     task.ID,
			UserID:     userID,
			AssignedBy: task.CreatorID,
			Progress:   0,
		}
	}
	utils.SuccessResponse(c, http.StatusOK, "OK", progress)
}
