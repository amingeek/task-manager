// backend/controllers/group_additional.go
package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
)

// GetGroupDetails - دریافت جزئیات یک گروه
func GetGroupDetails(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	var group models.Group
	if err := config.DB.Preload("Members.User").Preload("Creator").Preload("Tasks").
		First(&group, groupID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "گروه پیدا نشد")
		return
	}

	// بررسی عضویت کاربر در گروه
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND accepted = ?", groupID, userID, true).
		First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "شما عضو این گروه نیستید")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", group)
}

// UpdateGroup - بروزرسانی گروه
func UpdateGroup(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").
		First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌توانند گروه را بروزرسانی کنند")
		return
	}

	var group models.Group
	if err := config.DB.First(&group, groupID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "گروه پیدا نشد")
		return
	}

	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Description != "" {
		group.Description = req.Description
	}

	config.DB.Save(&group)
	utils.SuccessResponse(c, http.StatusOK, "گروه با موفقیت بروزرسانی شد", group)
}

// DeleteGroup - حذف گروه
func DeleteGroup(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").
		First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌توانند گروه را حذف کنند")
		return
	}

	var group models.Group
	if err := config.DB.First(&group, groupID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "گروه پیدا نشد")
		return
	}

	config.DB.Delete(&group)
	utils.SuccessResponse(c, http.StatusOK, "گروه با موفقیت حذف شد", nil)
}

// RemoveMember - حذف عضو از گروه
func RemoveMember(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")
	memberID := c.Param("user_id")

	// بررسی اینکه کاربر مدیر گروه است
	var adminMember models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").
		First(&adminMember).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌توانند اعضا را حذف کنند")
		return
	}

	// حذف عضو
	if err := config.DB.Where("group_id = ? AND user_id = ?", groupID, memberID).
		Delete(&models.GroupMember{}).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در حذف عضو")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "عضو با موفقیت از گروه حذف شد", nil)
}

// AcceptInvitation - پذیرفتن دعوت به گروه
func AcceptInvitation(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	// پیدا کردن دعوت
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "دعوت پیدا نشد")
		return
	}

	// اگر قبلاً پذیرفته شده
	if member.Accepted {
		utils.SuccessResponse(c, http.StatusOK, "شما قبلاً این دعوت را پذیرفته‌اید", member)
		return
	}

	// پذیرفتن دعوت
	member.Accepted = true
	if err := config.DB.Save(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در پذیرش دعوت")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "دعوت با موفقیت پذیرفته شد", member)
}

func GetPendingInvitations(c *gin.Context) {
	userID := c.GetUint("userID")

	var pendingMembers []models.GroupMember
	if err := config.DB.Where("user_id = ? AND accepted = ?", userID, false).
		Preload("Group").
		Preload("Group.Creator").
		Find(&pendingMembers).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load pending invitations")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", pendingMembers)
}
