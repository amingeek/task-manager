// backend/controllers/task_controller.go
package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
}

func GetTasks(c *gin.Context) {
	userID := c.GetUint("userID")
	var tasks []models.Task
	config.DB.Where("creator_id = ?", userID).Order("created_at DESC").Find(&tasks)
	utils.SuccessResponse(c, http.StatusOK, "OK", tasks)
}

func GetTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var task models.Task
	if err := config.DB.Where("id = ? AND creator_id = ?", taskID, userID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Database error")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OK", task)
}

func CreateTask(c *gin.Context) {
	userID := c.GetUint("userID")
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   userID,
		Status:      models.StatusPending,
	}
	config.DB.Create(&task)
	utils.SuccessResponse(c, http.StatusCreated, "Task created successfully", task)
}

func UpdateTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var task models.Task
	if err := config.DB.Where("id = ? AND creator_id = ?", taskID, userID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// Update fields if provided
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}

	config.DB.Save(&task)
	utils.SuccessResponse(c, http.StatusOK, "Task updated successfully", task)
}

func DeleteTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var task models.Task
	if err := config.DB.Where("id = ? AND creator_id = ?", taskID, userID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "Task not found")
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Database error")
		}
		return
	}

	config.DB.Delete(&task)
	utils.SuccessResponse(c, http.StatusOK, "Task deleted successfully", nil)
}
