// backend/controllers/analytics_controller.go

package controllers

import (
	"net/http"
	"time"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
)

type StreakData struct {
	CurrentStreak int `json:"current_streak"`
	LongestStreak int `json:"longest_streak"`
}

type AnalyticsSummary struct {
	TotalTasks      int64       `json:"total_tasks"`
	CompletedTasks  int64       `json:"completed_tasks"`
	PendingTasks    int64       `json:"pending_tasks"`
	InProgressTasks int64       `json:"in_progress_tasks"`
	CompletionRate  float64     `json:"completion_rate"`
	Streak          StreakData  `json:"streak"`
	GroupCount      int64       `json:"group_count"`
	GroupTaskCount  int64       `json:"group_task_count"`
}

// GetStreak - دریافت اطلاعات streak کاربر
func GetStreak(c *gin.Context) {
	userID := c.GetUint("userID")

	var streak models.Streak
	if err := config.DB.Where("user_id = ?", userID).First(&streak).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "اطلاعات streak پیدا نشد")
		return
	}

	// بررسی آخرین فعالیت
	if streak.LastActivityAt.Before(time.Now().AddDate(0, 0, -1)) {
		// بیش از یک روز گذشته است
		streak.CurrentStreak = 0
		config.DB.Save(&streak)
	}

	streakData := StreakData{
		CurrentStreak: streak.CurrentStreak,
		LongestStreak: streak.LongestStreak,
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", streakData)
}

// GetAnalyticsSummary - دریافت خلاصه تحلیل‌های کاربر
func GetAnalyticsSummary(c *gin.Context) {
	userID := c.GetUint("userID")

	var summary AnalyticsSummary

	// کل تسک‌های شخصی
	config.DB.Model(&models.Task{}).
		Where("creator_id = ? AND is_group_task = ?", userID, false).
		Count(&summary.TotalTasks)

	// تسک‌های کامل شده
	config.DB.Model(&models.Task{}).
		Where("creator_id = ? AND is_group_task = ? AND status = ?", userID, false, models.StatusCompleted).
		Count(&summary.CompletedTasks)

	// تسک‌های معلق
	config.DB.Model(&models.Task{}).
		Where("creator_id = ? AND is_group_task = ? AND status = ?", userID, false, models.StatusPending).
		Count(&summary.PendingTasks)

	// تسک‌های در حال انجام
	config.DB.Model(&models.Task{}).
		Where("creator_id = ? AND is_group_task = ? AND status = ?", userID, false, models.StatusInProgress).
		Count(&summary.InProgressTasks)

	// میزان تکمیل
	if summary.TotalTasks > 0 {
		summary.CompletionRate = float64(summary.CompletedTasks) / float64(summary.TotalTasks) * 100
	}

	// اطلاعات streak
	var streak models.Streak
	if err := config.DB.Where("user_id = ?", userID).First(&streak).Error; err == nil {
		summary.Streak.CurrentStreak = streak.CurrentStreak
		summary.Streak.LongestStreak = streak.LongestStreak
	}

	// تعداد گروه‌ها
	config.DB.Model(&models.GroupMember{}).
		Where("user_id = ? AND accepted = ?", userID, true).
		Count(&summary.GroupCount)

	// تعداد تسک‌های گروهی برای کاربر
	config.DB.Model(&models.TaskAssignment{}).
		Joins("JOIN tasks ON tasks.id = task_assignments.task_id").
		Where("task_assignments.user_id = ? AND tasks.is_group_task = ?", userID, true).
		Count(&summary.GroupTaskCount)

	utils.SuccessResponse(c, http.StatusOK, "OK", summary)
}
