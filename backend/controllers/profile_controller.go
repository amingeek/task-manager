// backend/controllers/profile_controller.go

package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"

	"github.com/gin-gonic/gin"
)

type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

// UpdateProfile - بروزرسانی پروفایل کاربر
func UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "کاربر پیدا نشد")
		return
	}

	// بروزرسانی فیلدها
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "خطا در بروزرسانی پروفایل")
		return
	}

	// پاک کردن کلمه عبور از response
	user.Password = ""

	utils.SuccessResponse(c, http.StatusOK, "پروفایل با موفقیت بروزرسانی شد", user)
}
