// backend/controllers/auth_controller.go

package controllers

import (
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"task-manager/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register - ثبت نام کاربر جدید
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, _ := utils.HashPassword(req.Password)
	user := models.User{Username: req.Username, Email: req.Email, Password: hashedPassword}
	config.DB.Create(&user)

	streak := models.Streak{UserID: user.ID, LastActivityAt: time.Now()}
	config.DB.Create(&streak)

	// توکن برای 24 ساعت
	token, _ := utils.GenerateToken(user.ID, user.Username, 24)

	utils.SuccessResponse(c, http.StatusCreated, "OK", gin.H{"token": token})
}

// Login - ورود کاربر
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if config.DB.Where("username = ?", req.Username).First(&user).Error != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid")
		return
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid")
		return
	}

	// توکن برای 24 ساعت
	token, _ := utils.GenerateToken(user.ID, user.Username, 24)

	utils.SuccessResponse(c, http.StatusOK, "OK", gin.H{"token": token})
}

// GetCurrentUser - دریافت اطلاعات کاربر فعلی
func GetCurrentUser(c *gin.Context) {
	userID := c.GetUint("userID")
	var user models.User
	config.DB.First(&user, userID)
	utils.SuccessResponse(c, http.StatusOK, "OK", user)
}
