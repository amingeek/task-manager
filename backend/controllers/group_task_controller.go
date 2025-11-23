// backend/controllers/group_task_controller.go

package controllers

import (
	"net/http"
	"strconv"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateGroupTaskRequest struct {
	Title        string     `json:"title" binding:"required"`
	Description  string     `json:"description"`
	DueDate      *time.Time `json:"due_date"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	UserIDs      []uint     `json:"user_ids"`                  // کاربران خاصی، خالی برای تمام اعضا
	RequireFiles bool       `json:"require_files"`             // آیا فایل آپلود الزامی است
	MaxFiles     int        `json:"max_files" binding:"min=1"` // حداکثر تعداد فایل
	AllowTypes   string     `json:"allow_types"`               // نوع‌های مجاز: pdf,image,video
}

type UpdateGroupTaskRequest struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	DueDate      *time.Time `json:"due_date"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	RequireFiles bool       `json:"require_files"`
	MaxFiles     int        `json:"max_files"`
	AllowTypes   string     `json:"allow_types"`
}

// CreateGroupTask - ایجاد تسک گروهی
func CreateGroupTask(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	var req CreateGroupTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ? AND accepted = ?", groupID, userID, "admin", true).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌تواند تسک ایجاد کنند")
		return
	}

	// تبدیل groupID به uint
	gid, err := strconv.ParseUint(groupID, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "شناسه گروه نامعتبر است")
		return
	}

	// تنظیم حداکثر فایل‌ها
	maxFiles := req.MaxFiles
	if maxFiles == 0 {
		maxFiles = 5
	}

	// تبدیل uint64 به uint
	groupIDUint := uint(gid)

	// ایجاد تسک گروهی
	task := models.Task{
		Title:        req.Title,
		Description:  req.Description,
		DueDate:      req.DueDate,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		CreatorID:    userID,
		IsGroupTask:  true,
		GroupID:      &groupIDUint,
		Status:       models.StatusPending,
		RequireFiles: req.RequireFiles,
		MaxFiles:     maxFiles,
		AllowTypes:   req.AllowTypes,
	}

	if err := config.DB.Create(&task).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در ایجاد تسک")
		return
	}

	// دریافت اعضای گروه برای اختصاص تسک
	var members []models.GroupMember
	if len(req.UserIDs) > 0 {
		config.DB.Where("group_id = ? AND user_id IN (?) AND accepted = ?", groupID, req.UserIDs, true).Find(&members)
	} else {
		config.DB.Where("group_id = ? AND accepted = ?", groupID, true).Find(&members)
	}

	// اختصاص تسک به اعضا
	for _, groupMember := range members {
		config.DB.Create(&models.TaskAssignment{
			TaskID: task.ID,
			UserID: groupMember.UserID,
		})

		// ایجاد رکورد پیشرفت برای هر عضو
		config.DB.Create(&models.GroupTaskProgress{
			TaskID:     task.ID,
			UserID:     groupMember.UserID,
			AssignedBy: userID,
		})
	}

	// دریافت تسک کامل
	var taskWithAssignments models.Task
	config.DB.Preload("TaskAssignments.User").
		Preload("GroupProgress").
		First(&taskWithAssignments, task.ID)

	utils.SuccessResponse(c, http.StatusCreated, "تسک گروهی با موفقیت ایجاد شد", taskWithAssignments)
}

// GetGroupTasks - دریافت تمام تسک‌های گروهی
func GetGroupTasks(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	// بررسی اینکه کاربر عضو گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND accepted = ?", groupID, userID, true).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "شما عضو این گروه نیستید")
		return
	}

	var tasks []models.Task
	config.DB.Where("group_id = ?", groupID).
		Preload("TaskAssignments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User")
		}).
		Preload("GroupProgress", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User")
		}).
		Preload("Creator").
		Preload("Files").
		Order("created_at DESC").
		Find(&tasks)

	utils.SuccessResponse(c, http.StatusOK, "OK", tasks)
}

// UpdateGroupTask - بروزرسانی تسک گروهی
func UpdateGroupTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskIDStr := c.Param("task_id")
	groupID := c.Param("id")

	var req UpdateGroupTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "شناسه تسک نامعتبر است")
		return
	}

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ? AND accepted = ?", groupID, userID, "admin", true).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌تواند تسک را بروزرسانی کنند")
		return
	}

	// بررسی وجود تسک
	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
		return
	}

	if task.CreatorID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط سازنده تسک می‌تواند آن را بروزرسانی کند")
		return
	}

	// بروزرسانی فیلدها
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.StartTime != nil {
		task.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		task.EndTime = req.EndTime
	}

	task.RequireFiles = req.RequireFiles
	if req.MaxFiles > 0 {
		task.MaxFiles = req.MaxFiles
	}
	if req.AllowTypes != "" {
		task.AllowTypes = req.AllowTypes
	}

	config.DB.Save(&task)

	// دریافت تسک کامل
	var updatedTask models.Task
	config.DB.Preload("TaskAssignments.User").
		Preload("GroupProgress").
		First(&updatedTask, taskID)

	utils.SuccessResponse(c, http.StatusOK, "تسک با موفقیت بروزرسانی شد", updatedTask)
}

// DeleteGroupTask - حذف تسک گروهی
func DeleteGroupTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskIDStr := c.Param("task_id")
	groupID := c.Param("id")

	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "شناسه تسک نامعتبر است")
		return
	}

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ? AND accepted = ?", groupID, userID, "admin", true).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌تواند تسک را حذف کنند")
		return
	}

	// بررسی وجود تسک
	var task models.Task
	if err := config.DB.First(&task, taskID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "تسک پیدا نشد")
		return
	}

	if task.CreatorID != userID {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط سازنده تسک می‌تواند آن را حذف کند")
		return
	}

	// حذف تسک و تمام وابستگی‌های آن
	config.DB.Delete(&task)

	utils.SuccessResponse(c, http.StatusOK, "تسک با موفقیت حذف شد", nil)
}

// UpdateTaskAssignment - بروزرسانی وضعیت انجام تسک توسط عضو
func UpdateTaskAssignment(c *gin.Context) {
	userID := c.GetUint("userID")
	assignmentID := c.Param("id")

	var assignment models.TaskAssignment
	if err := config.DB.Where("id = ? AND user_id = ?", assignmentID, userID).First(&assignment).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "انجام تسک پیدا نشد")
		return
	}

	assignment.Completed = !assignment.Completed
	config.DB.Save(&assignment)

	utils.SuccessResponse(c, http.StatusOK, "وضعیت تسک بروزرسانی شد", assignment)
}
