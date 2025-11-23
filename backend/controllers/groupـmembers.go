package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"task-manager/models"
	"task-manager/database"
)

// AddGroupMembers - اضافه کردن اعضا به گروه
func AddGroupMembers(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	var req struct {
		UserIds []uint `json:"user_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("📤 AddGroupMembers - GroupID: %d, UserIDs: %v\n", groupID, req.UserIds)

	// بررسی کن گروه موجود است
	var group models.Group
	if err := database.DB.First(&group, groupID).Error; err != nil {
		fmt.Printf("❌ Group not found: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	fmt.Printf("✅ Group found: %s\n", group.Name)

	// اضافه کردن هر عضو
	addedMembers := []models.GroupMember{}
	for _, userID := range req.UserIds {
		fmt.Printf("📌 Processing user: %d\n", userID)

		// بررسی کن کاربر موجود است
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			fmt.Printf("⚠️ User not found: %d - %v\n", userID, err)
			continue
		}

		// بررسی کن اعضای موجود
		var existingMember models.GroupMember
		existingCheck := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&existingMember)
		
		if existingCheck.Error == nil {
			fmt.Printf("⚠️ User %d already in group\n", userID)
			continue
		}

		// ایجاد عضو جدید
		newMember := models.GroupMember{
			GroupID: uint(groupID),
			UserID:  userID,
			Role:    "member",
		}

		if err := database.DB.Create(&newMember).Error; err != nil {
			fmt.Printf("❌ Error adding user %d: %v\n", userID, err)
			continue
		}

		fmt.Printf("✅ User %d added to group\n", userID)
		addedMembers = append(addedMembers, newMember)
	}

	fmt.Printf("✅ Total members added: %d\n", len(addedMembers))

	c.JSON(http.StatusOK, gin.H{
		"message": "Members added successfully",
		"added":   len(addedMembers),
	})
}

// RemoveGroupMember - حذف عضو از گروه
func RemoveGroupMember(c *gin.Context) {
	groupIDStr := c.Param("id")
	userIDStr := c.Param("userId")

	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	fmt.Printf("📤 RemoveGroupMember - GroupID: %d, UserID: %d\n", groupID, userID)

	// حذف عضو
	if err := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{}).Error; err != nil {
		fmt.Printf("❌ Error removing member: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	fmt.Printf("✅ Member removed successfully\n")

	c.JSON(http.StatusOK, gin.H{
		"message": "Member removed successfully",
	})
}

// GetGroupMembers - دریافت اعضای گروه
func GetGroupMembers(c *gin.Context) {
	groupIDStr := c.Param("id")
	groupID, err := strconv.ParseUint(groupIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	fmt.Printf("📥 GetGroupMembers - GroupID: %d\n", groupID)

	var members []models.GroupMember
	if err := database.DB.Where("group_id = ?", groupID).
		Preload("User").
		Find(&members).Error; err != nil {
		fmt.Printf("❌ Error fetching members: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
		return
	}

	fmt.Printf("✅ Found %d members\n", len(members))

	c.JSON(http.StatusOK, gin.H{
		"data": members,
	})
}
