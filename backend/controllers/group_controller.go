// backend/controllers/group_controller.go
package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
)

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	UserIDs     []uint `json:"user_ids"`
}

type AddMembersRequest struct {
	UserIDs []uint `json:"user_ids"`
}

func CreateGroup(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Validate required fields
	if req.Name == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Group name is required")
		return
	}

	// Create group
	group := models.Group{
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   userID,
	}

	if err := config.DB.Create(&group).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create group")
		return
	}

	// Add creator as admin
	creatorMember := models.GroupMember{
		GroupID:  group.ID,
		UserID:   userID,
		Accepted: true,
		Role:     "admin",
	}

	if err := config.DB.Create(&creatorMember).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to add group creator")
		return
	}

	// Add other members if provided
	if req.UserIDs != nil && len(req.UserIDs) > 0 {
		for _, memberID := range req.UserIDs {
			if memberID != userID { // Don't add creator again
				// Check if user exists
				var user models.User
				if err := config.DB.First(&user, memberID).Error; err != nil {
					continue // Skip if user doesn't exist
				}

				// Check if member already exists
				var existingMember models.GroupMember
				if err := config.DB.Where("group_id = ? AND user_id = ?", group.ID, memberID).First(&existingMember).Error; err != nil {
					// Member doesn't exist, create new
					newMember := models.GroupMember{
						GroupID:  group.ID,
						UserID:   memberID,
						Accepted: false,
						Role:     "member",
					}

					if err := config.DB.Create(&newMember).Error; err != nil {
						// Log error but continue with other members
						continue
					}
				}
			}
		}
	}

	// Load group with relations
	var createdGroup models.Group
	if err := config.DB.Preload("Members.User").Preload("Creator").First(&createdGroup, group.ID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Group created but failed to load details")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Group created successfully", createdGroup)
}

func GetUserGroups(c *gin.Context) {
	userID := c.GetUint("userID")
	var groups []models.Group

	if err := config.DB.Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ? AND group_members.accepted = ?", userID, true).
		Preload("Members.User").
		Preload("Creator").
		Preload("Tasks").
		Find(&groups).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load groups")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", groups)
}

func SearchGroups(c *gin.Context) {
	query := c.Query("q")
	var groups []models.Group

	if query != "" {
		if err := config.DB.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").
			Preload("Members.User").
			Preload("Creator").
			Find(&groups).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search groups")
			return
		}
	} else {
		if err := config.DB.Preload("Members.User").Preload("Creator").Find(&groups).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load groups")
			return
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", groups)
}

func AddGroupMembers(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")
	var req AddMembersRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user is group admin
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Only group admins can add members")
		return
	}

	// Add new members
	if req.UserIDs != nil && len(req.UserIDs) > 0 {
		for _, memberID := range req.UserIDs {
			// Check if user exists
			var user models.User
			if err := config.DB.First(&user, memberID).Error; err != nil {
				continue // Skip if user doesn't exist
			}

			var existingMember models.GroupMember
			if err := config.DB.Where("group_id = ? AND user_id = ?", groupID, memberID).First(&existingMember).Error; err != nil {
				// Member doesn't exist, create new
				newMember := models.GroupMember{
					GroupID:  member.GroupID,
					UserID:   memberID,
					Accepted: false,
					Role:     "member",
				}

				if err := config.DB.Create(&newMember).Error; err != nil {
					continue // Log error but continue with other members
				}
			}
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Members added successfully", nil)
}

// تابع جدید برای جستجوی کاربران
func SearchUsers(c *gin.Context) {
	query := c.Query("q")
	var users []models.User

	if query != "" {
		if err := config.DB.Where("username LIKE ? OR email LIKE ? OR full_name LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
			Find(&users).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search users")
			return
		}
	} else {
		if err := config.DB.Limit(20).Find(&users).Error; err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to load users")
			return
		}
	}

	// حذف اطلاعات حساس
	for i := range users {
		users[i].Password = ""
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", users)
}
