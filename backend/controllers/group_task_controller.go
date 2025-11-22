// backend/controllers/group_task_controller.go
package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateGroupTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	UserIDs     []uint     `json:"user_ids"` // Specific users to assign, empty for all group members
}

func CreateGroupTask(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")
	var req CreateGroupTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user is group member
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND accepted = ?", groupID, userID, true).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not a member of this group")
		return
	}

	// Create group task
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		CreatorID:   userID,
		IsGroupTask: true,
		GroupID:     &member.GroupID,
		Status:      models.StatusPending,
	}
	config.DB.Create(&task)

	// Get group members to assign
	var members []models.GroupMember
	if len(req.UserIDs) > 0 {
		config.DB.Where("group_id = ? AND user_id IN (?) AND accepted = ?", groupID, req.UserIDs, true).Find(&members)
	} else {
		config.DB.Where("group_id = ? AND accepted = ?", groupID, true).Find(&members)
	}

	// Assign task to members
	for _, groupMember := range members {
		config.DB.Create(&models.TaskAssignment{
			TaskID: task.ID,
			UserID: groupMember.UserID,
		})
	}

	// Load task with assignments
	var taskWithAssignments models.Task
	config.DB.Preload("TaskAssignments.User").First(&taskWithAssignments, task.ID)

	utils.SuccessResponse(c, http.StatusCreated, "Group task created successfully", taskWithAssignments)
}

func GetGroupTasks(c *gin.Context) {
	userID := c.GetUint("userID")
	groupID := c.Param("id")

	// Check if user is group member
	var member models.GroupMember
	if err := config.DB.Where("group_id = ? AND user_id = ? AND accepted = ?", groupID, userID, true).First(&member).Error; err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "You are not a member of this group")
		return
	}

	var tasks []models.Task
	config.DB.Where("group_id = ?", groupID).
		Preload("TaskAssignments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User")
		}).
		Preload("Creator").
		Order("created_at DESC").
		Find(&tasks)

	utils.SuccessResponse(c, http.StatusOK, "OK", tasks)
}

func UpdateTaskAssignment(c *gin.Context) {
	userID := c.GetUint("userID")
	assignmentID := c.Param("id")

	var assignment models.TaskAssignment
	if err := config.DB.Where("id = ? AND user_id = ?", assignmentID, userID).First(&assignment).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Task assignment not found")
		return
	}

	assignment.Completed = !assignment.Completed
	config.DB.Save(&assignment)

	utils.SuccessResponse(c, http.StatusOK, "Task status updated", assignment)
}
