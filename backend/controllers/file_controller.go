package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	userID := c.GetUint("userID")

	// دریافت فایل از فرم
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "فایل ارسال نشده است")
		return
	}

	taskIDStr := c.PostForm("task_id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "شناسه تسک نامعتبر است")
		return
	}

	notes := c.PostForm("notes")

	// بررسی وجود تسک
	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
		return
	}

	// بررسی دسترسی کاربر به تسک
	if task.IsGroupTask {
		// برای تسک گروهی، بررسی عضویت در گروه
		var member models.GroupMember
		if err := config.DB.Where("group_id = ? AND user_id = ?", task.GroupID, userID).First(&member).Error; err != nil {
			utils.ErrorResponse(c, http.StatusForbidden, "شما به این تسک دسترسی ندارید")
			return
		}
	} else {
		// برای تسک شخصی، بررسی مالکیت
		if task.CreatorID != userID {
			utils.ErrorResponse(c, http.StatusForbidden, "شما به این تسک دسترسی ندارید")
			return
		}
	}

	// ایجاد پوشه برای ذخیره فایل
	uploadDir := fmt.Sprintf("uploads/tasks/%d", taskID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در ایجاد پوشه")
		return
	}

	// تولید نام یکتا برای فایل
	fileName := fmt.Sprintf("%d_%d_%s", userID, time.Now().Unix(), file.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	// ذخیره فایل
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در ذخیره فایل")
		return
	}

	// ایجاد رکورد فایل در دیتابیس
	fileRecord := models.File{
		Filename: file.Filename,
		Filepath: filePath,
		FileSize: file.Size,
		MimeType: file.Header.Get("Content-Type"),
		UserID:   userID,
		TaskID:   uint(taskID),
	}

	config.DB.Create(&fileRecord)

	// اگر تسک گروهی است، آپدیت یادداشت در پیشرفت
	if task.IsGroupTask {
		var groupProgress models.GroupTaskProgress
		if err := config.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&groupProgress).Error; err != nil {
			// ایجاد رکورد جدید اگر وجود ندارد
			groupProgress = models.GroupTaskProgress{
				TaskID:     uint(taskID),
				UserID:     userID,
				AssignedBy: task.CreatorID,
				Notes:      notes,
			}
			config.DB.Create(&groupProgress)
		} else {
			groupProgress.Notes = notes
			config.DB.Save(&groupProgress)
		}
	}

	utils.SuccessResponse(c, http.StatusCreated, "فایل با موفقیت آپلود شد", fileRecord)
}

func GetTaskFiles(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
		return
	}

	// بررسی دسترسی
	if task.IsGroupTask {
		var member models.GroupMember
		if err := config.DB.Where("group_id = ? AND user_id = ?", task.GroupID, userID).First(&member).Error; err != nil {
			utils.ErrorResponse(c, http.StatusForbidden, "شما به این تسک دسترسی ندارید")
			return
		}
	} else {
		if task.CreatorID != userID {
			utils.ErrorResponse(c, http.StatusForbidden, "شما به این تسک دسترسی ندارید")
			return
		}
	}

	var files []models.File
	config.DB.Where("task_id = ?", taskID).Preload("User").Find(&files)

	utils.SuccessResponse(c, http.StatusOK, "OK", files)
}

func DownloadFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Preload("Task").First(&file, fileID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "فایل پیدا نشد")
		return
	}

	// بررسی دسترسی
	if file.Task.IsGroupTask {
		var member models.GroupMember
		if err := config.DB.Where("group_id = ? AND user_id = ?", file.Task.GroupID, userID).First(&member).Error; err != nil {
			utils.ErrorResponse(c, http.StatusForbidden, "شما به این فایل دسترسی ندارید")
			return
		}
	} else {
		if file.Task.CreatorID != userID && file.UserID != userID {
			utils.ErrorResponse(c, http.StatusForbidden, "شما به این فایل دسترسی ندارید")
			return
		}
	}

	// بررسی وجود فایل فیزیکی
	if _, err := os.Stat(file.Filepath); os.IsNotExist(err) {
		utils.ErrorResponse(c, http.StatusNotFound, "فایل روی سرور پیدا نشد")
		return
	}

	c.File(file.Filepath)
}

func DeleteFile(c *gin.Context) {
	userID := c.GetUint("userID")
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Preload("Task").First(&file, fileID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "فایل پیدا نشد")
		return
	}

	// فقط مالک فایل یا سازنده تسک می‌تواند حذف کند
	if file.UserID != userID && file.Task.CreatorID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "شما اجازه حذف این فایل را ندارید")
		return
	}

	// حذف فایل فیزیکی
	if _, err := os.Stat(file.Filepath); err == nil {
		os.Remove(file.Filepath)
	}

	// حذف رکورد دیتابیس
	config.DB.Delete(&file)

	utils.SuccessResponse(c, http.StatusOK, "فایل با موفقیت حذف شد", nil)
}
