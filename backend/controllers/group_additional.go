// backend/controllers/group_additional.go

package controllers

import (
	"net/http"
	"strconv"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
)

// GetGroups - دریافت تمام گروه‌های کاربر
func GetGroups(c *gin.Context) {
	userID := c.GetUint("userID")

	var groups []models.Group
	if err := config.DB.Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ? AND group_members.accepted = ?", userID, true).
		Preload("Members.User").
		Preload("Creator").
		Preload("Tasks").
		Find(&groups).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در دریافت گروه‌ها")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", groups)
}

// GetGroup - دریافت جزئیات یک گروه
func GetGroup(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	// بررسی اینکه کاربر عضو گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND accepted = ?", groupID, userID, true).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "شما عضو این گروه نیستید")
		return
	}

	var group models.Group
	if err := config.DB.Preload("Members.User").
		Preload("Creator").
		Preload("Tasks").
		First(&group, groupID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "گروه پیدا نشد")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", group)
}

// UpdateGroup - بروزرسانی اطلاعات گروه
func UpdateGroup(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌توانند گروه را بروزرسانی کنند")
		return
	}

	var group models.Group
	if err := config.DB.First(&group, groupID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "گروه پیدا نشد")
		return
	}

	// بروزرسانی فیلدها
	if req.Name != "" {
		group.Name = req.Name
	}
	if req.Description != "" {
		group.Description = req.Description
	}

	if err := config.DB.Save(&group).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در بروزرسانی گروه")
		return
	}

	// دریافت گروه کامل
	var updatedGroup models.Group
	config.DB.Preload("Members.User").
		Preload("Creator").
		First(&updatedGroup, groupID)

	utils.SuccessResponse(c, http.StatusOK, "گروه با موفقیت بروزرسانی شد", updatedGroup)
}

// DeleteGroup - حذف گروه
func DeleteGroup(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌توانند گروه را حذف کنند")
		return
	}

	var group models.Group
	if err := config.DB.First(&group, groupID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "گروه پیدا نشد")
		return
	}

	// حذف گروه و تمام وابستگی‌های آن
	if err := config.DB.Delete(&group).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در حذف گروه")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "گروه با موفقیت حذف شد", nil)
}

// AddMembers - اضافه کردن اعضای جدید به گروه
func AddMembers(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	var req AddMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌توانند اعضا را اضافه کنند")
		return
	}

	// اضافه کردن اعضای جدید
	for _, memberID := range req.UserIDs {
		// بررسی وجود کاربر
		var user models.User
		if err := config.DB.First(&user, memberID).Error; err != nil {
			continue
		}

		// بررسی اینکه عضو قبلاً اضافه نشده است
		var existingMember models.GroupMember
		if err := config.DB.Where("group_id = ? AND user_id = ?", groupID, memberID).First(&existingMember).Error; err != nil {
			// عضو جدید است
			newMember := models.GroupMember{
				GroupID:  convertStringToUint(groupID),
				UserID:   memberID,
				Accepted: false,
				Role:     "member",
			}
			config.DB.Create(&newMember)
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "اعضا با موفقیت اضافه شدند", nil)
}

// RemoveMember - حذف عضو از گروه
func RemoveMember(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")
	memberUserID := c.Param("user_id")

	// بررسی اینکه کاربر مدیر گروه است
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "فقط مدیران گروه می‌توانند اعضا را حذف کنند")
		return
	}

	// حذف عضو
	if err := config.DB.Where("group_id = ? AND user_id = ?", groupID, memberUserID).Delete(&models.GroupMember{}).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در حذف عضو")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "عضو با موفقیت حذف شد", nil)
}

// AcceptInvitation - پذیرش دعوت جوین گروه
func AcceptInvitation(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")
	memberUserID := c.Param("user_id")

	// بررسی اینکه کاربر درخواست کننده است
	if userID != convertStringToUint(memberUserID) {
		utils.ErrorResponse(c, http.StatusForbidden, "تنها خود کاربر می‌تواند دعوت را پذیرفته کند")
		return
	}

	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ?", groupID, memberUserID).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "دعوت پیدا نشد")
		return
	}

	member.Accepted = true
	if err := config.DB.Save(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در پذیرش دعوت")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "دعوت با موفقیت پذیرفته شد", nil)
}

// Helper function
func convertStringToUint(s string) uint {
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint(v)
}
