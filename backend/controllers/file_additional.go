// backend/controllers/file_additional.go

package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// ApproveGroupTaskFile - تایید فایل برای تسک گروهی توسط مدیر
func ApproveGroupTaskFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("file_id")
	groupID := c.Param("id")

	// دریافت فایل
	var file models.File
	if err := config.DB.Preload("Task").First(&file, fileID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "فایل پیدا نشد")
		return
	}

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌توانند فایل‌ها را تایید کنند")
		return
	}

	// بررسی اینکه تسک متعلق به گروه است
	var task models.Task
	if err := config.DB.First(&task, file.TaskID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
		return
	}

	if task.GroupID == nil || *task.GroupID != StringToUint(groupID) {
		utils.ErrorResponse(c, http.StatusForbidden, "این تسک متعلق به این گروه نیست")
		return
	}

	// تایید فایل
	file.Approved = true
	file.ApprovedBy = &userID
	now := time.Now()
	file.ApprovedAt = &now
	if err := config.DB.Save(&file).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در تایید فایل")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "فایل با موفقیت تایید شد", file)
}
